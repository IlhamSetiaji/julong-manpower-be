package route

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/handler"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App              *gin.Engine
	MPPPeriodHandler handler.IMPPPeriodHander
	AuthMiddleware   gin.HandlerFunc
}

func (c *RouteConfig) SetupRoutes() {
	c.App.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to Julong Manpower",
		})
	})
	c.SetupMPPPeriodRoutes()
}

func (c *RouteConfig) SetupMPPPeriodRoutes() {
	mppPeriod := c.App.Group("/api/mpp-periods")
	{
		mppPeriod.Use(c.AuthMiddleware)
		mppPeriod.GET("/", c.MPPPeriodHandler.FindAllPaginated)
		mppPeriod.GET("/:id", c.MPPPeriodHandler.FindById)
		mppPeriod.POST("/", c.MPPPeriodHandler.Create)
		mppPeriod.PUT("/", c.MPPPeriodHandler.Update)
		mppPeriod.DELETE("/:id", c.MPPPeriodHandler.Delete)
	}
}

func NewRouteConfig(app *gin.Engine, mppPeriodHandler handler.IMPPPeriodHander, authMiddleware gin.HandlerFunc) *RouteConfig {
	return &RouteConfig{
		App:              app,
		MPPPeriodHandler: mppPeriodHandler,
		AuthMiddleware:   authMiddleware,
	}
}
