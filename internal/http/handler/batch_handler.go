package handler

import (
	"net/http"
	"strconv"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/helper"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/middleware"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/usecase"
	"github.com/IlhamSetiaji/julong-manpower-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IBatchHandler interface {
	CreateBatchHeaderAndLines(c *gin.Context)
	GetBatchHeadersByStatus(c *gin.Context)
	FindByStatus(c *gin.Context)
	TriggerCreate(c *gin.Context)
	GetBatchedMPPlanningHeaders(c *gin.Context)
	FindById(c *gin.Context)
	FindDocumentByID(c *gin.Context)
	FindByNeedApproval(c *gin.Context)
	FindByCurrentDocumentDateAndStatus(c *gin.Context)
	UpdateStatusBatchHeader(c *gin.Context)
	GetCompletedBatchHeader(c *gin.Context)
	GetOrganizationsForBatchApproval(c *gin.Context)
}

type BatchHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	UseCase    usecase.IBatchUsecase
	Validate   *validator.Validate
	UserHelper helper.IUserHelper
}

func NewBatchHandler(log *logrus.Logger, viper *viper.Viper, useCase usecase.IBatchUsecase, validate *validator.Validate, userHelper helper.IUserHelper) IBatchHandler {
	return &BatchHandler{
		Log:        log,
		Viper:      viper,
		UseCase:    useCase,
		Validate:   validate,
		UserHelper: userHelper,
	}
}

func (h *BatchHandler) GetCompletedBatchHeader(c *gin.Context) {
	batch, err := h.UseCase.GetCompletedBatchHeader()

	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get completed batch header", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Completed batch header found", batch)
}

func (h *BatchHandler) GetBatchHeadersByStatus(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := c.Query("search")
	if search == "" {
		search = ""
	}

	createdAt := c.Query("created_at")
	if createdAt == "" {
		createdAt = "DESC"
	}

	sort := map[string]interface{}{
		"created_at": createdAt,
	}

	status := c.Query("status")

	if status == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", "Invalid request")
		return
	}

	approverType := c.Query("approver_type")
	if approverType == "" {
		approverType = string(entity.BatchHeaderApproverTypeCEO)
	}

	user, err := middleware.GetUser(c, h.Log)
	if err != nil {
		h.Log.Errorf("Error when getting user: %v", err)
		utils.ErrorResponse(c, 500, "error", err.Error())
		return
	}

	if user == nil {
		h.Log.Errorf("User not found")
		utils.ErrorResponse(c, 404, "error", "User not found")
		return
	}

	orgUUID, err := h.UserHelper.GetOrganizationID(user)
	if err != nil {
		h.Log.Errorf("Error when getting organization id: %v", err)
		utils.ErrorResponse(c, 500, "error", err.Error())
		return
	}

	batch, total, err := h.UseCase.GetBatchHeadersByStatusPaginated(entity.BatchHeaderApprovalStatus(status), approverType, orgUUID.String(), page, pageSize, search, sort)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get batch headers by status", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Batch headers found", gin.H{
		"batches": batch,
		"total":   total,
	})
}

func (h *BatchHandler) CreateBatchHeaderAndLines(c *gin.Context) {
	var req request.CreateBatchHeaderAndLinesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := middleware.GetUser(c, h.Log)
	if err != nil {
		h.Log.Errorf("Error when getting user: %v", err)
		utils.ErrorResponse(c, 500, "error", err.Error())
		return
	}
	if user == nil {
		h.Log.Errorf("User not found")
		utils.ErrorResponse(c, 404, "error", "User not found")
		return
	}

	orgUUID, err := h.UserHelper.GetOrganizationID(user)
	if err != nil {
		h.Log.Errorf("Error when getting organization id: %v", err)
		utils.ErrorResponse(c, 500, "error", err.Error())
		return
	}

	req.OrganizationID = orgUUID.String()

	batchHeader, err := h.UseCase.CreateBatchHeaderAndLines(&req)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create batch header and lines", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Batch header and lines created", batchHeader)
}

func (h *BatchHandler) GetOrganizationsForBatchApproval(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", "Invalid request")
		return
	}

	organizations, err := h.UseCase.GetOrganizationsForBatchApproval(id)

	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get organizations for batch approval", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Organizations for batch approval found", organizations)
}

func (h *BatchHandler) TriggerCreate(c *gin.Context) {
	approverType := c.Query("approver_type")
	if approverType == "" {
		approverType = string(entity.BatchHeaderApproverTypeCEO)
	}

	user, err := middleware.GetUser(c, h.Log)
	if err != nil {
		h.Log.Errorf("Error when getting user: %v", err)
		utils.ErrorResponse(c, 500, "error", err.Error())
		return
	}

	if user == nil {
		h.Log.Errorf("User not found")
		utils.ErrorResponse(c, 404, "error", "User not found")
		return
	}

	orgUUID, err := h.UserHelper.GetOrganizationID(user)
	if err != nil {
		h.Log.Errorf("Error when getting organization id: %v", err)
		utils.ErrorResponse(c, 500, "error", err.Error())
		return
	}

	batch, err := h.UseCase.TriggerCreate(approverType, orgUUID.String())
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to trigger create batch", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Batch created", batch)
}

func (h *BatchHandler) GetBatchedMPPlanningHeaders(c *gin.Context) {
	approverType := c.Query("approver_type")
	if approverType == "" {
		approverType = string(entity.BatchHeaderApproverTypeCEO)
	}

	user, err := middleware.GetUser(c, h.Log)
	if err != nil {
		h.Log.Errorf("Error when getting user: %v", err)
		utils.ErrorResponse(c, 500, "error", err.Error())
		return
	}

	if user == nil {
		h.Log.Errorf("User not found")
		utils.ErrorResponse(c, 404, "error", "User not found")
		return
	}

	orgUUID, err := h.UserHelper.GetOrganizationID(user)
	if err != nil {
		h.Log.Errorf("Error when getting organization id: %v", err)
		utils.ErrorResponse(c, 500, "error", err.Error())
		return
	}

	mpPlanningHeaders, err := h.UseCase.GetBatchedMPPlanningHeaders(approverType, orgUUID.String())
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get batched MP Planning headers", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Batched MP Planning headers found", mpPlanningHeaders)
}

func (h *BatchHandler) FindByStatus(c *gin.Context) {
	status := c.Param("status")
	approverType := c.Query("approver_type")
	if approverType == "" {
		approverType = string(entity.BatchHeaderApproverTypeCEO)
	}

	user, err := middleware.GetUser(c, h.Log)
	if err != nil {
		h.Log.Errorf("Error when getting user: %v", err)
		utils.ErrorResponse(c, 500, "error", err.Error())
		return
	}
	if user == nil {
		h.Log.Errorf("User not found")
		utils.ErrorResponse(c, 404, "error", "User not found")
		return
	}

	orgUUID, err := h.UserHelper.GetOrganizationID(user)
	if err != nil {
		h.Log.Errorf("Error when getting organization id: %v", err)
		utils.ErrorResponse(c, 500, "error", err.Error())
		return
	}

	approvalStatus := entity.BatchHeaderApprovalStatus(status)
	batch, err := h.UseCase.FindByStatus(approvalStatus, approverType, orgUUID.String())
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to find batch by status", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Batch found", batch)
}

func (h *BatchHandler) FindById(c *gin.Context) {
	id := c.Param("id")
	batch, err := h.UseCase.FindById(id)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to find batch by id", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Batch found", batch)
}

func (h *BatchHandler) FindDocumentByID(c *gin.Context) {
	id := c.Param("id")
	batch, err := h.UseCase.FindDocumentByID(id)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to find document by id", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Document found", batch)
}

func (h *BatchHandler) FindByNeedApproval(c *gin.Context) {
	approverType := c.Query("approver_type")
	if approverType == "" {
		approverType = string(entity.BatchHeaderApproverTypeCEO)
	}

	user, err := middleware.GetUser(c, h.Log)
	if err != nil {
		h.Log.Errorf("Error when getting user: %v", err)
		utils.ErrorResponse(c, 500, "error", err.Error())
		return
	}
	if user == nil {
		h.Log.Errorf("User not found")
		utils.ErrorResponse(c, 404, "error", "User not found")
		return
	}

	orgUUID, err := h.UserHelper.GetOrganizationID(user)
	if err != nil {
		h.Log.Errorf("Error when getting organization id: %v", err)
		utils.ErrorResponse(c, 500, "error", err.Error())
		return
	}
	batch, err := h.UseCase.FindByNeedApproval(approverType, orgUUID.String())
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to find batch by need approval", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Batch found", batch)
}

func (h *BatchHandler) FindByCurrentDocumentDateAndStatus(c *gin.Context) {
	status := c.Param("status")
	approvalStatus := entity.BatchHeaderApprovalStatus(status)
	batch, err := h.UseCase.FindByCurrentDocumentDateAndStatus(approvalStatus)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to find batch by current document date and status", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Batch found", batch)
}

func (h *BatchHandler) UpdateStatusBatchHeader(c *gin.Context) {
	var req request.UpdateStatusBatchHeaderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	batchExist, err := h.UseCase.FindById(req.ID)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to find batch by id", err.Error())
		return
	}

	if batchExist == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Batch not found", "Batch not found")
		return
	}

	if batchExist.Status == string(entity.BatchHeaderApprovalStatus(req.Status)) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Status already same", "Status already same")
		return
	}

	batchHeader, err := h.UseCase.UpdateStatusBatchHeader(&req)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update batch header status", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Batch header status updated", batchHeader)
}

func BatchHandlerFactory(log *logrus.Logger, viper *viper.Viper) IBatchHandler {
	validate := config.NewValidator(viper)
	userHelper := helper.NewUserHelper(log)
	useCase := usecase.BatchUsecaseFactory(viper, log)
	return NewBatchHandler(log, viper, useCase, validate, userHelper)
}
