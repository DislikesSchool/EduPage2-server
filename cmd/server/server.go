package main

import (
	"fmt"
	"net/http"
	"os"

	docs "github.com/DislikesSchool/EduPage2-server/docs"
	"github.com/DislikesSchool/EduPage2-server/edupage"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/robfig/cron/v3"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
)

type ClientData struct {
	CrJobId cron.EntryID
	Client  *edupage.EdupageClient
}

var clients = make(map[string]*ClientData)

var cr *cron.Cron

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
		EnableTracing:    false,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v", err)
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("EduPage2-server"),
		newrelic.ConfigLicense("eu01xx4155dca9b59d668e2cc9fc4e98FFFFNRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		fmt.Println("NewRelic initialization failed:", err)
	}

	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" && gin.Mode() == gin.ReleaseMode {
		fmt.Println("\033[0;31mERROR\033[0m: No JWT_SECRET_KEY environment variable found. Use the JWT_SECRET_KEY environment variable to set the secret key.")
		panic("No JWT_SECRET_KEY environment variable found")
	}

	cr = cron.New()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))
	router.Use(nrgin.Middleware(app))
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
	api.GET("/recipients", RecipientsHandler)
	api.POST("/message", SendMessageHandler)
	api.GET("/grades", ResultsHandler)

	router.POST("/icanteen", ICanteenHandler)

	router.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	router.StaticFile("/", "./cmd/server/web/index.html")
	router.StaticFile("/main.dart.js", "./cmd/server/web/main.dart.js")
	router.StaticFile("/flutter.js", "./cmd/server/web/flutter.js")
	router.StaticFile("/flutter_service_worker.js", "./cmd/server/web/flutter_service_worker.js")
	router.StaticFile("/manifest.json", "./cmd/server/web/manifest.json")
	router.StaticFile("/version.json", "./cmd/server/web/version.json")
	router.StaticFile("/favicon.png", "./cmd/server/web/favicon.png")
	router.StaticFS("/assets", http.Dir("./cmd/server/web/assets"))
	router.StaticFS("/canvaskit", http.Dir("./cmd/server/web/canvaskit"))
	router.StaticFS("/icons", http.Dir("./cmd/server/web/icons"))

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(302, "/")
	})

	cr.Start()
	router.Run()
}
