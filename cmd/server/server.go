package main

import (
	"fmt"
	"time"

	docs "github.com/DislikesSchool/EduPage2-server/docs"
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
