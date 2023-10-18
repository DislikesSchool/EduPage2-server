package main

import (
	"fmt"
	"os"

	docs "github.com/DislikesSchool/EduPage2-server/docs"
	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var clients = make(map[string]*edupage.EdupageClient)

// @title EduPage2 API
// @version 1.0
// @description This is the backend for the EduPage2 app.
// @BasePath /

// @SecurityDefinition Bearer
// @Description JWT authorization token
// @Type apiKey
// @In header
// @Name Authorization

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

	router.POST("/login", LoginHandler)
	router.GET("/validate-token", ValidateTokenHandler)
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api.GET("/timeline", TimelineHandler)
	api.GET("/timeline/recent", RecentTimelineHandler)
	api.GET("/timetable", TimetableHandler)
	api.GET("/timetable/recent", RecentTimetableHangler)
	api.GET("/subject/:id", SubjectHandler)
	api.GET("/teacher/:id", TeacherHandler)
	api.GET("/classroom/:id", ClassroomHandler)
	api.GET("/periods", PeriodsHandler)
	api.GET("/timelineitem/:id", TimelineItemHandler)

	router.POST("/icanteen", ICanteenHandler)

	router.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "test"})
	})

	router.Run()
}
