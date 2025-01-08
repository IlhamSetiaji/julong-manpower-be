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
	MPRequestHandler       handler.IMPRequestHandler
	BatchHandler           handler.IBatchHandler
	AuthMiddleware         gin.HandlerFunc
}

func (c *RouteConfig) SetupRoutes() {
	// c.App.GET("/", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"message": "Welcome to Julong Manpower",
	// 	})
	// })

	c.SetupAPIRoutes()
	// c.SetupMPPPeriodRoutes()
	// c.SetupJobPlafonRoutes()
	// c.SetupMPPlanningRoutes()
	// c.SetupRequestCategoryRoutes()
	// c.SetupMajorRoutes()
	// c.SetupMPRequestRoutes()
}

func (c *RouteConfig) SetupAPIRoutes() {
	apiRoute := c.App.Group("/api")
	{
		apiRoute.Use(c.AuthMiddleware)
		{
			apiRoute.Use(c.AuthMiddleware)
			apiRoute.GET("/mpp-periods", c.MPPPeriodHandler.FindAllPaginated)
			apiRoute.GET("/mpp-periods/current", c.MPPPeriodHandler.FindByCurrentDateAndStatus)
			apiRoute.GET("/mpp-periods/status", c.MPPPeriodHandler.FindByStatus)
			apiRoute.GET("/mpp-periods/update-status", c.MPPPeriodHandler.UpdateStatusByDate)
			apiRoute.GET("/mpp-periods/:id", c.MPPPeriodHandler.FindById)
			apiRoute.POST("/mpp-periods", c.MPPPeriodHandler.Create)
			apiRoute.PUT("/mpp-periods", c.MPPPeriodHandler.Update)
			apiRoute.DELETE("/mpp-periods/:id", c.MPPPeriodHandler.Delete)

			// job plafon
			apiRoute.GET("/job-plafons", c.JobPlafonHandler.FindAllPaginated)
			apiRoute.GET("/job-plafons/sync", c.JobPlafonHandler.SyncJobPlafon)
			apiRoute.GET("/job-plafons/:id", c.JobPlafonHandler.FindById)
			apiRoute.GET("/job-plafons/job/:job_id", c.JobPlafonHandler.FindByJobId)
			apiRoute.POST("/job-plafons", c.JobPlafonHandler.Create)
			apiRoute.PUT("/job-plafons", c.JobPlafonHandler.Update)
			apiRoute.DELETE("/job-plafons/:id", c.JobPlafonHandler.Delete)

			// mp plannings
			apiRoute.GET("/mp-plannings", c.MPPlanningHandler.FindAllHeadersPaginated)
			apiRoute.GET("/mp-plannings/something", c.MPPlanningHandler.FindHeaderBySomething)
			apiRoute.GET("/mp-plannings/get-something", c.MPPlanningHandler.GetHeadersBySomething)
			apiRoute.GET("/mp-plannings/status-period", c.MPPlanningHandler.FindAllHeadersByStatusAndMPPeriodID)
			apiRoute.GET("/mp-plannings/completed", c.MPPlanningHandler.GetHeadersByMPPeriodCompleted)
			apiRoute.GET("/mp-plannings/total-histories", c.MPPlanningHandler.CountTotalApprovalHistoryByStatus)
			apiRoute.GET("/mp-plannings/document-number", c.MPPlanningHandler.GenerateDocumentNumber)
			apiRoute.GET("/mp-plannings/requestor", c.MPPlanningHandler.FindAllHeadersByRequestorIDPaginated)
			apiRoute.GET("/mp-plannings/batch", c.MPPlanningHandler.FindAllHeadersForBatchPaginated)
			apiRoute.GET("/mp-plannings/count", c.MPPlanningHandler.CountMPPlanningHeaderByMPPPeriodIDAndApproverType)
			apiRoute.GET("/mp-plannings/approver-type", c.MPPlanningHandler.FindAllHeadersGroupedApproverPaginated)
			apiRoute.GET("/mp-plannings/jobs/:header_id", c.MPPlanningHandler.FindJobsByHeaderID)
			apiRoute.GET("/mp-plannings/locations/:header_id", c.MPPlanningHandler.FindOrganizationLocationsByHeaderID)
			apiRoute.GET("/mp-plannings/mpp-period/:mpp_period_id", c.MPPlanningHandler.FindHeaderByMPPPeriodId)
			apiRoute.GET("/mp-plannings/approval-attachments/:approval_history_id", c.MPPlanningHandler.GetPlanningApprovalHistoryAttachmentsByApprovalHistoryId)
			apiRoute.GET("/mp-plannings/approval-histories/:header_id", c.MPPlanningHandler.GetPlanningApprovalHistoryByHeaderId)
			apiRoute.GET("/mp-plannings/:id", c.MPPlanningHandler.FindById)
			apiRoute.POST("/mp-plannings", c.MPPlanningHandler.Create)
			apiRoute.PUT("/mp-plannings", c.MPPlanningHandler.Update)
			apiRoute.PUT("/mp-plannings/update-status", c.MPPlanningHandler.UpdateStatusMPPPlanningHeader)
			apiRoute.DELETE("/mp-plannings/:id", c.MPPlanningHandler.Delete)

			// mp planning lines
			apiRoute.GET("/mp-plannings/lines/find/:id", c.MPPlanningHandler.FindLineById)
			apiRoute.GET("/mp-plannings/lines/:header_id", c.MPPlanningHandler.FindAllLinesByHeaderIdPaginated)
			apiRoute.POST("/mp-plannings/lines/store", c.MPPlanningHandler.CreateLine)
			apiRoute.PUT("/mp-plannings/lines/update", c.MPPlanningHandler.UpdateLine)
			apiRoute.PUT("/mp-plannings/lines/reject-partial", c.MPPlanningHandler.RejectStatusPartialMPPlanningHeader)
			apiRoute.PUT("/mp-plannings/lines/reject-partial-pt", c.MPPlanningHandler.RejectStatusPartialMPPlanningHeaderUsingPT)
			apiRoute.DELETE("/mp-plannings/lines/delete/:id", c.MPPlanningHandler.DeleteLine)
			apiRoute.POST("/mp-plannings/lines/batch/store", c.MPPlanningHandler.CreateOrUpdateBatchLineMPPlanningLines)

			// request categories
			apiRoute.GET("/request-categories", c.RequestCategoryHandler.FindAll)
			apiRoute.GET("/request-categories/is-replacement", c.RequestCategoryHandler.GetByIsReplacement)
			apiRoute.GET("/request-categories/:id", c.RequestCategoryHandler.FindById)

			// majors
			apiRoute.GET("/majors", c.MajorHandler.FindAll)
			apiRoute.GET("/majors/education-level", c.MajorHandler.GetMajorsByEducationLevel)
			apiRoute.GET("/majors/:id", c.MajorHandler.FindById)

			// mp requests
			apiRoute.GET("/mp-requests", c.MPRequestHandler.FindAllPaginated)
			apiRoute.GET("/mp-requests/total-histories", c.MPRequestHandler.CountTotalApprovalHistoryByStatus)
			apiRoute.GET("/mp-requests/document-number", c.MPRequestHandler.GenerateDocumentNumber)
			apiRoute.GET("/mp-requests/approval-histories", c.MPRequestHandler.GetRequestApprovalHistoryByHeaderId)
			apiRoute.GET("/mp-requests/:id", c.MPRequestHandler.FindByID)
			apiRoute.POST("/mp-requests", c.MPRequestHandler.Create)
			apiRoute.PUT("/mp-requests/status", c.MPRequestHandler.UpdateStatusMPRequestHeader)
			apiRoute.PUT("/mp-requests", c.MPRequestHandler.Update)
			apiRoute.DELETE("/mp-requests/:id", c.MPRequestHandler.Delete)

			// batch
			apiRoute.POST("/batch/create", c.BatchHandler.CreateBatchHeaderAndLines)
			apiRoute.GET("/batch/trigger-create", c.BatchHandler.TriggerCreate)
			apiRoute.GET("/batch/batched-list", c.BatchHandler.GetBatchedMPPlanningHeaders)
			apiRoute.GET("/batch/completed", c.BatchHandler.GetCompletedBatchHeader)
			apiRoute.GET("/batch/need-approval", c.BatchHandler.FindByNeedApproval)
			apiRoute.GET("/batch/organizations/:id", c.BatchHandler.GetOrganizationsForBatchApproval)
			apiRoute.GET("/batch/find-by-status/:status", c.BatchHandler.FindByStatus)
			apiRoute.GET("/batch/find-document/:id", c.BatchHandler.FindDocumentByID)
			apiRoute.GET("/batch/current-status/:status", c.BatchHandler.FindByCurrentDocumentDateAndStatus)
			apiRoute.GET("/batch/:id", c.BatchHandler.FindById)
			apiRoute.PUT("/batch/update-status", c.BatchHandler.UpdateStatusBatchHeader)
		}
	}
}

func NewRouteConfig(app *gin.Engine, viper *viper.Viper, log *logrus.Logger) *RouteConfig {
	// factory handlers
	mppPeriodHandler := handler.MPPPeriodHandlerFactory(log, viper)
	jobPlafonHandler := handler.JobPlafonHandlerFactory(log, viper)
	mpPlanningHandler := handler.MPPlanningHandlerFactory(log, viper)
	requestCategoryHandler := handler.RequestCategoryHandlerFactory(log, viper)
	majorHandler := handler.MajorHandlerFactory(log, viper)
	mpRequestHandler := handler.MPRequestHandlerFactory(log, viper)
	batchHandler := handler.BatchHandlerFactory(log, viper)

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
		MPRequestHandler:       mpRequestHandler,
		BatchHandler:           batchHandler,
	}
}
