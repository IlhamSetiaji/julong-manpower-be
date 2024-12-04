package main

import (
	"net/http"
	"strconv"

	"github.com/IlhamSetiaji/go-rabbitmq-utils/rabbitmq"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/handler"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/middleware"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/route"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func main() {
	viper := config.NewViper()
	log := config.NewLogrus(viper)

	err := rabbitmq.InitializeConnection(viper.GetString("rabbitmq.url"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.CloseConnection()

	log.Info("RabbitMQ connection established")

	app := gin.Default()
	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("App-Name", viper.GetString("app.name"))
	})

	store := cookie.NewStore([]byte(viper.GetString("web.cookie.secret")))
	app.Use(sessions.Sessions(viper.GetString("web.session.name"), store))

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

	// factory handlers
	mppPeriodHandler := handler.MPPPeriodHandlerFactory(log, viper)
	jobPlafonHandler := handler.JobPlafonHandlerFactory(log, viper)

	// facroty middleware
	authMiddleware := middleware.NewAuth(viper)

	// setup routes
	routeConfig := route.NewRouteConfig(app, mppPeriodHandler, authMiddleware, jobPlafonHandler)
	routeConfig.SetupRoutes()

	// run server
	webPort := strconv.Itoa(viper.GetInt("web.port"))
	err = app.Run(":" + webPort)
	if err != nil {
		log.Panicf("Failed to start server: %v", err)
	}
}

func shouldExcludeFromCSRF(path string) bool {
	return len(path) >= 4 && path[:4] == "/api"
}
