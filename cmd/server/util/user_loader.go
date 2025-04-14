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
	"github.com/meilisearch/meilisearch-go"
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

	// Set up periodic message indexing for search if enabled
	if ShouldSearch && ShouldStore && Cr != nil {
		// Schedule message indexing to run once a day at 3 AM to minimize impact
		_, err := Cr.AddFunc("0 3 * * *", updateAllUsersMessageIndex)
		if err != nil {
			ErrorLogger.Printf("Failed to schedule message indexing job: %v", err)
		} else {
			InfoLogger.Println("Scheduled daily message index updates at 3 AM")
		}
	}

	updateAllUsersMessageIndex()
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

func updateAllUsersMessageIndex() {
	if !ShouldSearch || !ShouldStore {
		return
	}

	var users []dbmodel.User
	result := Db.Where("store_messages = ?", true).
		Where("last_online > ?", time.Now().AddDate(0, 0, -7)).
		Find(&users)

	if result.Error != nil {
		ErrorLogger.Printf("Failed to fetch users for message indexing: %v", result.Error)
		return
	}

	InfoLogger.Printf("Starting scheduled message indexing for %d users", len(users))

	// Process users sequentially to avoid overwhelming the server
	for _, user := range users {
		LoadMessagesAndUpdateMeilisearch(&user)
		// Small delay between users to reduce load
		time.Sleep(5 * time.Second)
	}

	InfoLogger.Println("Completed scheduled message indexing")
}

func LoadMessagesAndUpdateMeilisearch(user *dbmodel.User) {
	if !ShouldSearch {
		return
	}

	// Generate a unique owner ID for this user
	ownerID := fmt.Sprintf("%s@%s", user.Username, user.Server)

	// Get client for this user
	clientKey := user.Server + user.Username
	clientData, exists := Clients[clientKey]
	if !exists {
		InfoLogger.Printf("Cannot load messages for %s: client not found", ownerID)
		return
	}

	InfoLogger.Printf("Loading messages for %s into search index", ownerID)

	// Get messages from the past year
	endDate := time.Now()
	startDate := endDate.AddDate(-1, 0, 0) // One year ago

	timeline, err := clientData.Client.GetTimeline(startDate, endDate)
	if err != nil {
		ErrorLogger.Printf("Failed to load messages for %s: %v", ownerID, err)
		return
	}

	// Check if we have any messages
	if len(timeline.Items) == 0 {
		InfoLogger.Printf("No messages found for %s", ownerID)
		return
	}

	// Prepare documents for Meilisearch
	var messagesToIndex []map[string]interface{}

	// Get existing message IDs to avoid duplicates
	existingIDs, err := getExistingMessageIDs(ownerID)
	if err != nil {
		ErrorLogger.Printf("Failed to get existing message IDs: %v", err)
		// Continue anyway, might result in duplicates but that's better than no indexing
	}

	// Process messages
	for id, item := range timeline.Items {
		// Skip if not a message type
		if item.Type != "message" && item.Type != "sprava" && !strings.Contains(strings.ToLower(item.Type), "message") {
			continue
		}

		// Skip if already indexed
		if existingIDs[id] {
			continue
		}

		// Create a document with only the necessary fields for searching
		doc := map[string]interface{}{
			"id":          id,
			"ownerid":     ownerID,        // Important for filtering
			"text":        item.Text,      // Message content for searching
			"sender_id":   item.Owner,     // Sender ID for filtering
			"receiver_id": item.User,      // Receiver ID for filtering
			"timestamp":   item.Timestamp, // For sorting results
		}

		messagesToIndex = append(messagesToIndex, doc)
	}

	// If we have messages to index, send them to Meilisearch
	if len(messagesToIndex) > 0 {
		InfoLogger.Printf("Indexing %d new messages for %s", len(messagesToIndex), ownerID)

		// Add documents to Meilisearch index
		_, err := MeiliIndex.AddDocuments(messagesToIndex)
		if err != nil {
			ErrorLogger.Printf("Failed to index messages for %s: %v", ownerID, err)
			return
		}

		InfoLogger.Printf("Successfully indexed messages for %s", ownerID)
	} else {
		InfoLogger.Printf("No new messages to index for %s", ownerID)
	}
}

// Helper function to get existing message IDs for this user to avoid duplicates
func getExistingMessageIDs(ownerID string) (map[string]bool, error) {
	existingIDs := make(map[string]bool)

	// Search for documents with this ownerID
	searchRequest := &meilisearch.SearchRequest{
		Filter: fmt.Sprintf("ownerid = \"%s\"", ownerID),
		Limit:  10000, // Adjust based on expected message volume
	}

	searchResponse, err := MeiliIndex.Search("", searchRequest)
	if err != nil {
		return existingIDs, err
	}

	// Extract IDs from search results
	for _, hit := range searchResponse.Hits {
		if doc, ok := hit.(map[string]interface{}); ok {
			if id, exists := doc["id"].(string); exists {
				existingIDs[id] = true
			}
		}
	}

	return existingIDs, nil
}
