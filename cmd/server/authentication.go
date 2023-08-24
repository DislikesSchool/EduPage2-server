package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func getSecretKey() []byte {
	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" {
		panic("JWT_SECRET_KEY environment variable is not set")
	}
	return []byte(key)
}

func generateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(time.Hour)

	claims := jwt.MapClaims{
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

func verifyJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return getSecretKey(), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
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

	c.JSON(http.StatusOK, gin.H{
		"error":   "",
		"success": true,
		"name":    h.EdupageData.User.UserRow.Firstname + " " + h.EdupageData.User.UserRow.Lastname,
	})
}
