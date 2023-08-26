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

func generateJWT(userID string, username string) (string, error) {
	expirationTime := time.Now().Add(time.Hour)

	claims := jwt.MapClaims{
		"userID":   userID,
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

		c.Next()
	}
}

func getUserIDAndUsername(c *gin.Context) (string, string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", "", errors.New("missing Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return getSecretKey(), nil
	})
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", "", errors.New("invalid user ID in token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", "", errors.New("invalid user ID in token claims")
	}

	return userID, username, nil
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
// @Accept json
// @Accept multipart/form-data
// @Accept x-www-form-urlencoded
// @Param login body LoginRequestUsernamePassword false "Login using username and password"
// @Param loginServer body LoginRequestUsernamePasswordServer false "Login using username, password and server"
// @Produce json
// @Success 200 {object} LoginSuccessResponse
// @Failure 400 {object} LoginBadRequestResponse
// @Failure 401 {object} LoginUnauthorizedResponse
// @Failure 500 {object} LoginInternalErrorResponse
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	var loginData LoginData
	if err := c.Bind(&loginData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if loginData.Username == "" || loginData.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username and Password are required"})
		return
	}

	var h edupage.EdupageClient
	var err error
	if loginData.Server == "" {
		h, err = edupage.LoginAuto(loginData.Username, loginData.Password)
	} else {
		h, err = edupage.Login(loginData.Server, loginData.Username, loginData.Password)
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	err = h.LoadUser()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	userID := h.EdupageData.User.UserRow.UserID
	username := loginData.Username

	token, err := generateJWT(userID, loginData.Username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	clients[userID+username] = &h
	c.JSON(http.StatusOK, gin.H{
		"error":   "",
		"success": true,
		"name":    h.EdupageData.User.UserRow.Firstname + " " + h.EdupageData.User.UserRow.Lastname,
		"token":   token,
	})
}
