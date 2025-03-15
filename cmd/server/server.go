package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/DislikesSchool/EduPage2-server/config"
	docs "github.com/DislikesSchool/EduPage2-server/docs"
	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ClientData struct {
	CrJobId cron.EntryID
	Client  *edupage.EdupageClient
}

var clients = make(map[string]*ClientData)

var ctx = context.Background()

var cr *cron.Cron

var rdb redis.Client

var shouldCache = config.AppConfig.Redis.Enabled

// @title EduPage2 API
// @version 1.1.0
// @description This is the backend for the EduPage2 app.
// @BasePath /

// @SecurityDefinition Bearer
// @Description JWT authorization token
// @Type apiKey
// @In header
// @Name Authorization

func main() {
	key := config.AppConfig.JWT.Secret
	if key == "" && gin.Mode() == gin.ReleaseMode {
		fmt.Println("\033[0;31mERROR\033[0m: No JWT_SECRET_KEY environment variable found. Use the JWT_SECRET_KEY environment variable to set the secret key.")
		panic("No JWT_SECRET_KEY environment variable found")
	}

	cr = cron.New()

	if config.AppConfig.Redis.Enabled {
		rdb = *redis.NewClient(&redis.Options{
			Addr:     config.AppConfig.Redis.Address,
			Username: config.AppConfig.Redis.Username,
			Password: config.AppConfig.Redis.Password,
			DB:       config.AppConfig.Redis.DB,
		})
	}

	if config.AppConfig.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/test"}}),
		gin.Recovery(),
	)
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))
	api := router.Group("/api")
	api.Use(authMiddleware())
	docs.SwaggerInfo.BasePath = "/"

	router.POST("/login", LoginHandler)
	router.GET("/validate-token", ValidateTokenHandler)
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/qrlogin", QRLoginHandler)
	router.POST("/qrlogin/:code", FinishQRLoginHandler)

	api.GET("/timeline", TimelineHandler)
	api.GET("/timeline/recent", RecentTimelineHandler)
	api.GET("/timetable", TimetableHandler)
	api.GET("/timetable/recent", RecentTimetableHangler)
	api.GET("/subject/:id", SubjectHandler)
	api.GET("/teacher/:id", TeacherHandler)
	api.GET("/classroom/:id", ClassroomHandler)
	api.GET("/periods", PeriodsHandler)
	api.GET("/timelineitem/:id", TimelineItemHandler)
	api.GET("/recipients", RecipientsHandler)
	api.POST("/message", SendMessageHandler)
	api.GET("/grades", ResultsHandler)

	ic := router.Group("/icanteen")
	ic.POST("/login", ICanteenLoginHandler)
	ic.POST("/month", ICanteenMonthHandler)
	ic.POST("/change", ICanteenChangeOrderHandler)

	// For compatibility with 1.0.x
	router.POST("/icanteen", ICanteenHandler)
	router.POST("/icanteen-test", ICanteenTestHandler)

	router.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	router.StaticFile("/", "./cmd/server/web/index.html")
	router.StaticFile("/main.dart.js", "./cmd/server/web/main.dart.js")
	router.StaticFile("/flutter.js", "./cmd/server/web/flutter.js")
	router.StaticFile("/flutter_bootstrap.js", "./cmd/server/web/flutter.js")
	router.StaticFile("/flutter_service_worker.js", "./cmd/server/web/flutter_service_worker.js")
	router.StaticFile("/manifest.json", "./cmd/server/web/manifest.json")
	router.StaticFile("/version.json", "./cmd/server/web/version.json")
	router.StaticFile("/favicon.png", "./cmd/server/web/favicon.png")
	router.StaticFS("/assets", http.Dir("./cmd/server/web/assets"))
	router.StaticFS("/canvaskit", http.Dir("./cmd/server/web/canvaskit"))
	router.StaticFS("/icons", http.Dir("./cmd/server/web/icons"))

	router.StaticFile("/.well-known/assetlinks.json", "./cmd/server/.well-known/assetlinks.json")

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "/")
	})

	cr.Start()

	port := config.AppConfig.Server.Port
	if port == "" {
		port = "8080"
	}
	host := config.AppConfig.Server.Host
	if host == "" {
		host = "0.0.0.0"
	}
	router.Run(host + ":" + port)
}
