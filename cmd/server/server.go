package main

import (
	"fmt"
	"net/http"

	"github.com/DislikesSchool/EduPage2-server/cmd/server/crypto"
	"github.com/DislikesSchool/EduPage2-server/cmd/server/dbmodel"
	"github.com/DislikesSchool/EduPage2-server/cmd/server/routes"
	"github.com/DislikesSchool/EduPage2-server/cmd/server/util"
	"github.com/DislikesSchool/EduPage2-server/config"
	docs "github.com/DislikesSchool/EduPage2-server/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// @title EduPage2 API
// @version 1.2.0
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

	util.Cr = cron.New()

	if util.ShouldCache {
		util.Rdb = *redis.NewClient(&redis.Options{
			Addr:     config.AppConfig.Redis.Address,
			Username: config.AppConfig.Redis.Username,
			Password: config.AppConfig.Redis.Password,
			DB:       config.AppConfig.Redis.DB,
		})
	}

	if config.AppConfig.Encryption.Enabled {
		if config.AppConfig.Encryption.Key == "" {
			fmt.Println("\033[0;31mERROR\033[0m: No encryption key found. Use a command like openssl rand -base64 32 to generate a key.")
			panic("No encryption key found")
		} else if config.AppConfig.Encryption.Key == "YourDefaultEncryptionKey" {
			fmt.Println("\033[0;31mERROR\033[0m: Using default encryption key. Please change it for security reasons.")
			panic("Using default encryption key")
		}
		if err := crypto.InitCrypto(config.AppConfig.Encryption.Key); err != nil {
			fmt.Printf("\033[0;31mERROR\033[0m: %v\n", err)
			panic("Failed to initialize encryption")
		}
	}

	if util.ShouldStore {
		dsn := config.AppConfig.Database.DSN
		var err error
		var dialector gorm.Dialector
		switch config.AppConfig.Database.Driver {
		case "sqlite":
			dialector = sqlite.Open(dsn)
		case "mysql":
			dialector = mysql.Open(dsn)
		case "postgres":
			dialector = postgres.Open(dsn)
		case "sqlserver":
			dialector = sqlserver.Open(dsn)
		default:
			panic("Unknown database driver")
		}
		util.Db, err = gorm.Open(dialector, &gorm.Config{})
		if err != nil {
			panic(err)
		}

		util.Db.AutoMigrate(&dbmodel.User{})
	}

	if config.AppConfig.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	util.InitLogging()
	defer util.CloseLogger()

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

	api.GET("/timeline", routes.TimelineHandler)
	api.GET("/timeline/recent", routes.RecentTimelineHandler)
	api.GET("/timetable", routes.TimetableHandler)
	api.GET("/timetable/recent", routes.RecentTimetableHandler)
	api.GET("/subject/:id", routes.SubjectHandler)
	api.GET("/teacher/:id", routes.TeacherHandler)
	api.GET("/classroom/:id", routes.ClassroomHandler)
	api.GET("/periods", routes.PeriodsHandler)
	api.GET("/timelineitem/:id", routes.TimelineItemHandler)
	api.GET("/recipients", routes.RecipientsHandler)
	api.POST("/message", routes.SendMessageHandler)
	api.GET("/grades", routes.ResultsHandler)

	api.GET("/version", routes.ServerVersion)
	api.GET("/capabilities", routes.ServerCapabilities)

	ic := router.Group("/icanteen")
	ic.POST("/login", routes.ICanteenLoginHandler)
	ic.POST("/month", routes.ICanteenMonthHandler)
	ic.POST("/change", routes.ICanteenChangeOrderHandler)

	// For compatibility with 1.0.x
	router.POST("/icanteen", routes.ICanteenHandler)
	router.POST("/icanteen-test", routes.ICanteenTestHandler)

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
		c.File("./cmd/server/web/index.html")
	})

	util.Cr.Start()

	if util.ShouldStore {
		util.InfoLogger.Println("Starting to load stored users...")
		util.LoadStoredUsers()
	}

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
