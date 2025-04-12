package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/DislikesSchool/EduPage2-server/cmd/server/crypto"
	"github.com/DislikesSchool/EduPage2-server/cmd/server/dbmodel"
	"github.com/DislikesSchool/EduPage2-server/cmd/server/util"
	"github.com/DislikesSchool/EduPage2-server/config"
	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func getSecretKey() []byte {
	key := config.AppConfig.JWT.Secret
	if key == "" {
		key = "development-secret-key"
	}
	return []byte(key)
}

func generateJWT(server string, username string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 6)

	claims := jwt.MapClaims{
		"server":   server,
		"username": username,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(getSecretKey())
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is invalid"})
			return
		}

		tokenString := authHeaderParts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return getSecretKey(), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		client, dataStorage, err := clientFromContext(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("client", client)
		c.Set("dataStorage", dataStorage)

		c.Next()
	}
}

func getClaims(c *gin.Context) (jwt.MapClaims, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getSecretKey(), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Server   string `json:"server" binding:"omitempty" required:"false"`
	Token    string `json:"token" binding:"omitempty" required:"false"`
}

// LoginHandler godoc
// @Summary Login to your Edupage account
// @Schemes
// @Description Logs in to your Edupage account using the provided credentials.
// @Tags auth
// @Accept multipart/form-data
// @Accept x-www-form-urlencoded
// @Consumes application/x-www-form-urlencoded
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param server formData string false "Server"
// @Param storage query string false "Data Storage config. If omitted, no data will be stored."
// @Produce json
// @Success 200 {object} apimodel.LoginSuccessResponse
// @Failure 400 {object} apimodel.LoginBadRequestResponse
// @Failure 401 {object} apimodel.LoginUnauthorizedResponse
// @Failure 500 {object} apimodel.LoginInternalErrorResponse
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username field is missing"})
		return
	}

	password := c.PostForm("password")
	if password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "password field is missing"})
		return
	}

	server := c.PostForm("server")

	if username == "" || password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username and Password are required"})
		return
	}

	var cred edupage.Credentials
	var err error
	if server == "" {
		server = "login1"
	}

	u := util.Clients[server+username]

	if u != nil {
		dataStorage := u.DataStorage

		if util.ShouldStore {
			// Try to get the 'storage' query param (stringified json), parse it and set it to dataStorage. If not present, leave it as is.
			if storageParam := c.Query("storage"); storageParam != "" {
				if err := json.Unmarshal([]byte(storageParam), dataStorage); err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid storage configuration"})
					return
				}
			}

			if dataStorage.Enabled && dataStorage.Credentials {
				userModel := &dbmodel.User{}
				result := util.Db.First(userModel, "username = ?", username)

				if result.Error == nil {
					userModel.LastOnline = time.Now()
					userModel.StoreMessages = dataStorage.Messages
					userModel.StoreTimeline = dataStorage.Timeline

					if err := util.Db.Save(userModel).Error; err != nil {
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user data: " + err.Error()})
						return
					}
				} else {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
					return
				}
			}
		}

		passwordCorrect := edupage.CheckPasswordHash(password, u.Client.Credentials.PasswordHash)
		if passwordCorrect {
			user, err := u.Client.GetUser(false)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error":   err.Error(),
					"success": false,
				})
				return
			}
			token, err := generateJWT(server, username)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"error":   "",
				"success": true,
				"name":    user.UserRow.Firstname + " " + user.UserRow.Lastname,
				"token":   token,
			})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "",
				"success": false,
			})
			return
		}
	}

	cred, err = edupage.Login(username, password, server)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	schoolId := strings.Split(cred.Server, ".")[0]
	if config.AppConfig.Schools.IsBlacklist {
		if slices.Contains(config.AppConfig.Schools.Whitelist, schoolId) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "This school is blacklisted",
				"success": false,
			})
			return
		}
	} else {
		if !slices.Contains(config.AppConfig.Schools.Whitelist, schoolId) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "This school is not whitelisted",
				"success": false,
			})
			return
		}
	}

	var h *edupage.EdupageClient
	h, err = edupage.CreateClient(cred)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	user, err := h.GetUser(false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	token, err := generateJWT(server, username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dataStorage := &util.DataStorageConfig{
		Enabled:     false,
		Credentials: false,
		Messages:    false,
		Timeline:    false,
	}

	if util.ShouldStore {
		// Try to get the 'storage' query param (stringified json), parse it and set it to dataStorage. If not present, leave it as is.
		if storageParam := c.Query("storage"); storageParam != "" {
			if err := json.Unmarshal([]byte(storageParam), dataStorage); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid storage configuration"})
				return
			}
		}

		if dataStorage.Enabled && dataStorage.Credentials {
			userModel := &dbmodel.User{}
			result := util.Db.First(userModel, "username = ?", username)

			var passwordToStore string
			if config.AppConfig.Encryption.Enabled {
				encryptedPwd, err := crypto.Encrypt(password)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password: " + err.Error()})
					return
				}
				passwordToStore = encryptedPwd
			} else {
				passwordToStore = password
			}

			if result.Error == nil {
				if userModel.Password != passwordToStore {
					userModel.Password = passwordToStore
				}
				if userModel.Server != server {
					userModel.Server = server
				}

				userModel.LastOnline = time.Now()
				userModel.StoreMessages = dataStorage.Messages
				userModel.StoreTimeline = dataStorage.Timeline

				if err := util.Db.Save(userModel).Error; err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user data: " + err.Error()})
					return
				}
			} else if result.Error == gorm.ErrRecordNotFound {
				newUser := &dbmodel.User{
					Username:      username,
					Password:      passwordToStore,
					Server:        server,
					LastOnline:    time.Now(),
					StoreMessages: dataStorage.Messages,
					StoreTimeline: dataStorage.Timeline,
				}

				if err := util.Db.Create(newUser).Error; err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
					return
				}
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}
		}
	}

	util.Clients[server+username] = &util.ClientData{
		Client:      h,
		DataStorage: dataStorage,
	}

	// The cron is kinda broken when run from the go test command
	if os.Getenv("CI") == "" {
		jobId, err := util.Cr.AddFunc("@every 10m", func() {
			fmt.Println("Pinging", username, strings.Split(cred.Server, ".")[0])
			success, err := h.PingSession()
			if err != nil || !success {
				fmt.Println("session ping failed")
				util.Cr.Remove(util.Clients[server+username].CrJobId)
				util.Clients[server+username] = nil
			}
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		util.Clients[server+username].CrJobId = jobId
	}
	c.JSON(http.StatusOK, gin.H{
		"error":   "",
		"success": true,
		"name":    user.UserRow.Firstname + " " + user.UserRow.Lastname,
		"token":   token,
	})
}

func clientFromContext(c *gin.Context) (*edupage.EdupageClient, *util.DataStorageConfig, error) {
	claims, err := getClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return &edupage.EdupageClient{}, &util.DataStorageConfig{}, err
	}
	server := claims["server"].(string)
	username := claims["username"].(string)

	client, ok := util.Clients[server+username]
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "client not found"})
		return &edupage.EdupageClient{}, &util.DataStorageConfig{}, err
	}

	return client.Client, client.DataStorage, nil
}

// ValidateTokenHandler godoc
// @Summary Validate your token
// @Schemes
// @Description Validates your token and returns a 200 OK if it's valid.
// @Tags auth
// @Param token header string true "JWT token"
// @Produce json
// @Success 200 {object} apimodel.ValidateTokenSuccessResponse
// @Failure 401 {object} apimodel.ValidateTokenUnauthorizedResponse
// @Router /validate-token [get]
func ValidateTokenHandler(c *gin.Context) {
	claims, err := getClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	server := claims["server"].(string)
	username := claims["username"].(string)
	exp := claims["exp"].(float64)

	h, ok := util.Clients[server+username]
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "client not found"})
		return
	}
	user, err := h.Client.GetUser(false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   "",
		"name":    user.UserRow.Firstname + " " + user.UserRow.Lastname,
		"success": true,
		"expires": exp,
	})
}

var (
	sseclients = make(map[string]chan QRLoginData)
	mu         sync.Mutex
)

type QRLoginData struct {
	Code     string `json:"code"`
	Username string `json:"username"`
	Password string `json:"password"`
	Endpoint string `json:"endpoint"`
	Server   string `json:"server"`
}

// QRLoginHandler godoc
// @Summary Log in using a QR code
// @Schemes
// @Description Logs in using a QR code. This route uses Server-Sent Events (SSE).
// @Tags auth
// @Router /qrlogin [get]
func QRLoginHandler(c *gin.Context) {
	code := generateCode(8)
	dataChannel := make(chan QRLoginData, 1)

	mu.Lock()
	sseclients[code] = dataChannel
	mu.Unlock()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	sentCode := false
	c.Stream(func(w io.Writer) bool {
		if !sentCode {
			c.SSEvent("code", code)
			sentCode = true
		}
		select {
		case data := <-dataChannel:
			c.SSEvent("data", data)
			return false
		default:
			return true
		}
	})

	mu.Lock()
	delete(sseclients, code)
	mu.Unlock()

	close(dataChannel)

	c.Status(http.StatusOK)
}

// FinishQRLoginHandler godoc
// @Summary Finish QR login
// @Schemes
// @Description Finishes QR login by sending the login data to the client that initiated the SSE channel.
// @Tags auth
// @Param code path string true "Code"
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Param endpoint formData string true "Endpoint"
// @Param server formData string true "Server"
// @Router /qrlogin/:code [post]
func FinishQRLoginHandler(c *gin.Context) {
	code := c.Param("code")
	username := c.PostForm("username")
	password := c.PostForm("password")
	endpoint := c.PostForm("endpoint")
	server := c.PostForm("server")

	mu.Lock()
	dataChannel, ok := sseclients[code]
	mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid code"})
		return
	}

	dataChannel <- QRLoginData{
		Code:     code,
		Username: username,
		Password: password,
		Endpoint: endpoint,
		Server:   server,
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login data sent"})
}

func generateCode(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
