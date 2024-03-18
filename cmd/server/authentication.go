package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func getSecretKey() []byte {
	key := os.Getenv("JWT_SECRET_KEY")
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

		client, err := clientFromContext(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("client", client)

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

	u := clients[server+username]

	if u != nil {
		passwordCorrect := edupage.CheckPasswordHash(password, u.Credentials.PasswordHash)
		if passwordCorrect {
			user, err := u.GetUser(false)
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

	clients[server+username] = h
	cr.AddFunc("@every 10m", func() {
		fmt.Println("Pinging", username, server)
		success, err := h.PingSession()
		if err != nil {
			fmt.Println("session ping failed")
			clients[server+username] = nil
		} else if !success {
			fmt.Println("session ping failed")
			clients[server+username] = nil
		}
	})
	c.JSON(http.StatusOK, gin.H{
		"error":   "",
		"success": true,
		"name":    user.UserRow.Firstname + " " + user.UserRow.Lastname,
		"token":   token,
	})
}

func clientFromContext(c *gin.Context) (*edupage.EdupageClient, error) {
	claims, err := getClaims(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return &edupage.EdupageClient{}, err
	}
	server := claims["server"].(string)
	username := claims["username"].(string)

	client, ok := clients[server+username]
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "client not found"})
		return &edupage.EdupageClient{}, err
	}

	return client, nil
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

	h, ok := clients[server+username]
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "client not found"})
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

	c.JSON(http.StatusOK, gin.H{
		"error":   "",
		"name":    user.UserRow.Firstname + " " + user.UserRow.Lastname,
		"success": true,
		"expires": exp,
	})
}
