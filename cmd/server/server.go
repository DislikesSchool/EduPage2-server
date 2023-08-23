package main

import (
	"fmt"
	"time"

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
		// Detect if the content is form-urlencoded or json, and use the appropriate binding function
		var loginData LoginData
		if err := c.ShouldBind(&loginData); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
	}))

	router.Run()
}
