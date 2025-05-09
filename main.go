package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	// "github.com/IlhamSetiaji/go-rabbitmq-utils/rabbitmq"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/route"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/scheduler"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/rabbitmq"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	csrf "github.com/utrack/gin-csrf"
)

func main() {
	viper := config.NewViper()
	log := config.NewLogrus(viper)

	go rabbitmq.InitConsumer(viper, log)
	go rabbitmq.InitProducer(viper, log)

	// err := rabbitmq.InitializeConnection(viper.GetString("rabbitmq.url"))
	// if err != nil {
	// 	log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	// }
	// defer rabbitmq.CloseConnection()

	// log.Info("RabbitMQ connection established")

	app := gin.Default()
	app.Static("/storage", "./storage")
	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("App-Name", viper.GetString("app.name"))
	})

	store := cookie.NewStore([]byte(viper.GetString("web.cookie.secret")))
	app.Use(sessions.Sessions(viper.GetString("web.session.name"), store))

	// setup CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(viper.GetString("frontend.urls"), ","), // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// setup custom csrf middleware
	app.Use(func(c *gin.Context) {
		if !shouldExcludeFromCSRF(c.Request.URL.Path) {
			csrf.Middleware(csrf.Options{
				Secret: viper.GetString("web.csrf_secret"),
				ErrorFunc: func(c *gin.Context) {
					c.String(http.StatusForbidden, "CSRF token mismatch")
					c.Abort()
				},
			})(c)
		}
		c.Next()
	})

	// setup routes
	routeConfig := route.NewRouteConfig(app, viper, log)
	routeConfig.SetupRoutes()

	// setup cron & scheduler
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	sch := cron.New(cron.WithLocation(jakartaTime))

	schedulerFactory := scheduler.MPPPeriodSchedulerFactory(log)
	_, err := sch.AddFunc("1 0 * * *", func() {
		err := schedulerFactory.UpdateStatusToOpenByDate()
		if err != nil {
			log.Errorf("Failed to update status to open by date: %v", err)
		}
		err = schedulerFactory.UpdateStatusToCloseByDate()
		if err != nil {
			log.Errorf("Failed to update status to close by date: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("failed to add cron job: %v", err)
	}

	sch.Start()
	log.Infof("Started cron job")
	defer sch.Stop()

	// run server
	webPort := strconv.Itoa(viper.GetInt("web.port"))
	err = app.Run(":" + webPort)
	if err != nil {
		log.Panicf("Failed to start server: %v", err)
	}

	select {}
}

func shouldExcludeFromCSRF(path string) bool {
	return len(path) >= 4 && path[:4] == "/api"
}
