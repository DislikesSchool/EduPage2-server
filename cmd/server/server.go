package main

import (
	"fmt"
	"os"
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

var clients = make(map[string]*edupage.EdupageClient)

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

	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" && gin.Mode() == gin.ReleaseMode {
		fmt.Println("\033[0;31mERROR\033[0m: No JWT_SECRET_KEY environment variable found. Use the JWT_SECRET_KEY environment variable to set the secret key.")
		panic("No JWT_SECRET_KEY environment variable found")
	}

	router := gin.Default()
	api := router.Group("/api")
	api.Use(authMiddleware())
	docs.SwaggerInfo.BasePath = "/"
	router.Use(sentrygin.New(sentrygin.Options{}))
	store := persistence.NewInMemoryStore(time.Second)

	router.POST("/login", cache.CachePage(store, time.Hour, LoginHandler))
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api.GET("/timeline/recent", cache.CachePage(store, time.Minute, RecentTimelineHandler))

	router.Run()
}
