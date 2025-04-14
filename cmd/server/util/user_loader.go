package util

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DislikesSchool/EduPage2-server/cmd/server/crypto"
	"github.com/DislikesSchool/EduPage2-server/cmd/server/dbmodel"
	"github.com/DislikesSchool/EduPage2-server/config"
	"github.com/DislikesSchool/EduPage2-server/edupage"
)

// LoadStoredUsers loads all stored users from the database and authenticates them
func LoadStoredUsers() {
	if !ShouldStore {
		return
	}

	var users []dbmodel.User
	result := Db.Find(&users)
	if result.Error != nil {
		LogUserSession("database load", "all", "all", result.Error)
		return
	}

	InfoLogger.Printf("Loading %d stored users...", len(users))

	// Use a goroutine pool to load users concurrently but with limits
	workerCount := 5 // Number of concurrent workers
	jobs := make(chan dbmodel.User, len(users))
	results := make(chan struct {
		username string
		server   string
		success  bool
		err      error
	}, len(users))

	// Start workers
	for w := 1; w <= workerCount; w++ {
		go func(id int) {
			for user := range jobs {
				err := loadAndAuthenticateUserWithRetry(&user, 3)
				results <- struct {
					username string
					server   string
					success  bool
					err      error
				}{
					username: user.Username,
					server:   user.Server,
					success:  err == nil,
					err:      err,
				}
			}
		}(w)
	}

	// Send jobs to workers
	activeUserCount := 0
	for _, user := range users {
		// Skip users that haven't been online in over 7 days
		if user.LastOnline.Before(time.Now().AddDate(0, 0, -7)) {
			continue
		}
		jobs <- user
		activeUserCount++
	}
	close(jobs)

	// Collect results
	successCount := 0
	for i := 0; i < activeUserCount; i++ {
		result := <-results
		if result.success {
			successCount++
			LogUserSession("auto-login", result.username, result.server, nil)
		} else {
			LogUserSession("auto-login", result.username, result.server, result.err)
		}
	}

	InfoLogger.Printf("Successfully loaded %d/%d active users", successCount, activeUserCount)
}

// loadAndAuthenticateUser loads a single user and authenticates them
func loadAndAuthenticateUser(user *dbmodel.User) error {
	// Decrypt password if encryption is enabled
	password := user.Password
	if config.AppConfig.Encryption.Enabled {
		var err error
		password, err = crypto.Decrypt(password)
		if err != nil {
			return fmt.Errorf("failed to decrypt password: %w", err)
		}
	}

	// Try to authenticate with Edupage
	cred, err := edupage.Login(user.Username, password, user.Server)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Create the Edupage client
	client, err := edupage.CreateClient(cred)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Create data storage config
	dataStorage := &DataStorageConfig{
		Enabled:     true,
		Credentials: true,
		Messages:    user.StoreMessages,
		Timeline:    user.StoreTimeline,
	}

	// Add the client to the clients map
	clientKey := user.Server + user.Username
	Clients[clientKey] = &ClientData{
		Client:      client,
		DataStorage: dataStorage,
	}

	// Set up periodic session pinging
	if Cr != nil && os.Getenv("CI") == "" {
		jobId, err := Cr.AddFunc("@every 10m", func() {
			fmt.Println("Pinging", user.Username, user.Server)
			success, err := client.PingSession()
			if err != nil || !success {
				fmt.Println("Session ping failed for", user.Username)
				Cr.Remove(Clients[clientKey].CrJobId)
				delete(Clients, clientKey)
			}
		})
		if err != nil {
			return fmt.Errorf("failed to add cron job: %w", err)
		}
		Clients[clientKey].CrJobId = jobId
	}

	return nil
}

func loadAndAuthenticateUserWithRetry(user *dbmodel.User, maxRetries int) error {
	var err error
	backoff := 2 * time.Second

	if user.Password == "" || user.Username == "" {
		return fmt.Errorf("user %s has no password or username", user.Username)
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		err = loadAndAuthenticateUser(user)
		if err == nil {
			return nil
		}

		// If it's not a temporary error, don't retry
		if !isTemporaryError(err) {
			return err
		}

		fmt.Printf("Temporary error loading user %s (attempt %d/%d): %v. Retrying in %v...\n",
			user.Username, attempt+1, maxRetries, err, backoff)

		time.Sleep(backoff)
		backoff *= 2 // Exponential backoff
		if backoff > 30*time.Second {
			backoff = 30 * time.Second // Cap at 30 seconds
		}
	}

	return fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
}

// Check if an error is likely temporary and worth retrying
func isTemporaryError(err error) bool {
	errString := err.Error()
	// Customize this list based on your experience with temporary errors
	temporaryErrors := []string{
		"connection reset by peer",
		"timeout",
		"too many requests",
		"server is busy",
		"connection refused",
		"no route to host",
		"temporary failure",
	}

	for _, temp := range temporaryErrors {
		if strings.Contains(strings.ToLower(errString), temp) {
			return true
		}
	}

	return false
}
