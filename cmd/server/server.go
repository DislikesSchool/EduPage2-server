package main

import (
	"fmt"
	"net/http"
	"time"

	docs "github.com/DislikesSchool/EduPage2-server/docs"
	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title EduPage2 API
// @version 1.0
// @description This is the backend for the EduPage2 app.
// @BasePath /

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

func main() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://9f278010f63fd37cc43bee40a6d69aa6@o4504950085976064.ingest.sentry.io/4505752992743424",
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v", err)
	}

	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/"
	router.Use(sentrygin.New(sentrygin.Options{}))
	store := persistence.NewInMemoryStore(time.Second)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/login", cache.CachePage(store, time.Minute, LoginHandler))

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run()
}
