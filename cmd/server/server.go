package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Server   string `json:"server"`
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
	router.Use(sentrygin.New(sentrygin.Options{}))
	store := persistence.NewInMemoryStore(time.Second)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/login", cache.CachePage(store, time.Minute, func(c *gin.Context) {
		var loginData LoginData
		if err := c.Bind(&loginData); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if loginData.Username == "" || loginData.Password == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username and Password are required"})
			return
		}

		var h edupage.Handle
		var err error
		if loginData.Server == "" {
			h, err = edupage.LoginAuto(loginData.Username, loginData.Password)
		} else {
			h, err = edupage.Login(loginData.Server, loginData.Username, loginData.Password)
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		data := h.RefreshUser()

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"name":    data.UserRow.Firstname + " " + data.UserRow.Lastname,
		})
	}))

	router.Run()
}
