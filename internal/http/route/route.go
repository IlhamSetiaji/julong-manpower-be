package route

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/handler"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App              *gin.Engine
	MPPPeriodHandler handler.IMPPPeriodHander
	JobPlafonHandler handler.IJobPlafonHandler
	AuthMiddleware   gin.HandlerFunc
}

func (c *RouteConfig) SetupRoutes() {
	c.App.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to Julong Manpower",
		})
	})
	c.SetupMPPPeriodRoutes()
	c.SetupJobPlafonRoutes()
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

func (c *RouteConfig) SetupJobPlafonRoutes() {
	jobPlafon := c.App.Group("/api/job-plafons")
	{
		jobPlafon.Use(c.AuthMiddleware)
		jobPlafon.GET("/", c.JobPlafonHandler.FindAllPaginated)
		jobPlafon.GET("/:id", c.JobPlafonHandler.FindById)
		jobPlafon.GET("/job/:job_id", c.JobPlafonHandler.FindByJobId)
		jobPlafon.POST("/", c.JobPlafonHandler.Create)
		jobPlafon.PUT("/", c.JobPlafonHandler.Update)
		jobPlafon.DELETE("/:id", c.JobPlafonHandler.Delete)
	}
}

func NewRouteConfig(app *gin.Engine, mppPeriodHandler handler.IMPPPeriodHander, authMiddleware gin.HandlerFunc, jobHandler handler.IJobPlafonHandler) *RouteConfig {
	return &RouteConfig{
		App:              app,
		MPPPeriodHandler: mppPeriodHandler,
		AuthMiddleware:   authMiddleware,
		JobPlafonHandler: jobHandler,
	}
}
