package route

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/handler"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App              *gin.Engine
	MPPPeriodHandler handler.IMPPPeriodHander
}

func (c *RouteConfig) SetupRoutes() {
	c.App.Group("/api")
	{
		c.SetupMPPPeriodRoutes()
	}
}

func (c *RouteConfig) SetupMPPPeriodRoutes() {
	mppPeriod := c.App.Group("/mpp-period")
	{
		mppPeriod.GET("/", c.MPPPeriodHandler.FindAllPaginated)
		mppPeriod.GET("/:id", c.MPPPeriodHandler.FindById)
		mppPeriod.POST("/", c.MPPPeriodHandler.Create)
		mppPeriod.PUT("/:id", c.MPPPeriodHandler.Update)
		mppPeriod.DELETE("/:id", c.MPPPeriodHandler.Delete)
	}
}

func NewRouteConfig(app *gin.Engine, mppPeriodHandler handler.IMPPPeriodHander) *RouteConfig {
	return &RouteConfig{
		App:              app,
		MPPPeriodHandler: mppPeriodHandler,
	}
}
