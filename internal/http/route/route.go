package route

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/handler"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	App                    *gin.Engine
	Log                    *logrus.Logger
	Viper                  *viper.Viper
	MPPPeriodHandler       handler.IMPPPeriodHander
	JobPlafonHandler       handler.IJobPlafonHandler
	MPPlanningHandler      handler.IMPPlanningHandler
	RequestCategoryHandler handler.IRequestCategoryHandler
	MajorHandler           handler.IMajorHandler
	AuthMiddleware         gin.HandlerFunc
}

func (c *RouteConfig) SetupRoutes() {
	c.App.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to Julong Manpower",
		})
	})
	c.SetupMPPPeriodRoutes()
	c.SetupJobPlafonRoutes()
	c.SetupMPPlanningRoutes()
	c.SetupRequestCategoryRoutes()
	c.SetupMajorRoutes()
}

func (c *RouteConfig) SetupMPPPeriodRoutes() {
	mppPeriod := c.App.Group("/api/mpp-periods")
	{
		mppPeriod.Use(c.AuthMiddleware)
		mppPeriod.GET("/", c.MPPPeriodHandler.FindAllPaginated)
		mppPeriod.GET("/current", c.MPPPeriodHandler.FindByCurrentDateAndStatus)
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

func (c *RouteConfig) SetupMPPlanningRoutes() {
	mpPlanning := c.App.Group("/api/mp-plannings")
	{
		mpPlanning.Use(c.AuthMiddleware)
		mpPlanning.GET("/", c.MPPlanningHandler.FindAllHeadersPaginated)
		mpPlanning.GET("/document-number", c.MPPlanningHandler.GenerateDocumentNumber)
		mpPlanning.GET("/requestor", c.MPPlanningHandler.FindAllHeadersByRequestorIDPaginated)
		mpPlanning.GET("/mpp-period/:mpp_period_id", c.MPPlanningHandler.FindHeaderByMPPPeriodId)
		mpPlanning.GET("/:id", c.MPPlanningHandler.FindById)
		mpPlanning.POST("/", c.MPPlanningHandler.Create)
		mpPlanning.PUT("/", c.MPPlanningHandler.Update)
		mpPlanning.DELETE("/:id", c.MPPlanningHandler.Delete)

		mpPlanning.GET("/lines/find/:id", c.MPPlanningHandler.FindLineById)
		mpPlanning.GET("/lines/:header_id", c.MPPlanningHandler.FindAllLinesByHeaderIdPaginated)
		mpPlanning.POST("/lines/store", c.MPPlanningHandler.CreateLine)
		mpPlanning.PUT("/lines/update", c.MPPlanningHandler.UpdateLine)
		mpPlanning.DELETE("/lines/delete/:id", c.MPPlanningHandler.DeleteLine)
		mpPlanning.POST("/lines/batch/store", c.MPPlanningHandler.CreateOrUpdateBatchLineMPPlanningLines)
	}
}

func (c *RouteConfig) SetupRequestCategoryRoutes() {
	requestCategory := c.App.Group("/api/request-categories")
	{
		requestCategory.Use(c.AuthMiddleware)
		requestCategory.GET("/", c.RequestCategoryHandler.FindAll)
		requestCategory.GET("/:id", c.RequestCategoryHandler.FindById)
	}
}

func (c *RouteConfig) SetupMajorRoutes() {
	major := c.App.Group("/api/majors")
	{
		major.Use(c.AuthMiddleware)
		major.GET("/", c.MajorHandler.FindAll)
		major.GET("/:id", c.MajorHandler.FindById)
	}
}

func NewRouteConfig(app *gin.Engine, viper *viper.Viper, log *logrus.Logger) *RouteConfig {
	// factory handlers
	mppPeriodHandler := handler.MPPPeriodHandlerFactory(log, viper)
	jobPlafonHandler := handler.JobPlafonHandlerFactory(log, viper)
	mpPlanningHandler := handler.MPPlanningHandlerFactory(log, viper)
	requestCategoryHandler := handler.RequestCategoryHandlerFactory(log, viper)
	majorHandler := handler.MajorHandlerFactory(log, viper)

	// facroty middleware
	authMiddleware := middleware.NewAuth(viper)
	return &RouteConfig{
		App:                    app,
		MPPPeriodHandler:       mppPeriodHandler,
		AuthMiddleware:         authMiddleware,
		JobPlafonHandler:       jobPlafonHandler,
		MPPlanningHandler:      mpPlanningHandler,
		RequestCategoryHandler: requestCategoryHandler,
		MajorHandler:           majorHandler,
	}
}
