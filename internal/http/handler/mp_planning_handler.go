package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/helper"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/middleware"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/usecase"
	"github.com/IlhamSetiaji/julong-manpower-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IMPPlanningHandler interface {
	FindAllHeadersPaginated(ctx *gin.Context)
	FindAllHeadersByRequestorIDPaginated(ctx *gin.Context)
	GenerateDocumentNumber(ctx *gin.Context)
	FindById(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	UpdateStatusMPPPlanningHeader(ctx *gin.Context)
	GetPlanningApprovalHistoryByHeaderId(ctx *gin.Context)
	FindHeaderByMPPPeriodId(ctx *gin.Context)
	FindAllLinesByHeaderIdPaginated(ctx *gin.Context)
	FindLineById(ctx *gin.Context)
	CreateLine(ctx *gin.Context)
	UpdateLine(ctx *gin.Context)
	DeleteLine(ctx *gin.Context)
	CreateOrUpdateBatchLineMPPlanningLines(ctx *gin.Context)
}

type MPPlanningHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	UseCase    usecase.IMPPlanningUseCase
	Validate   *validator.Validate
	UserHelper helper.IUserHelper
}

func NewMPPlanningHandler(log *logrus.Logger, viper *viper.Viper, useCase usecase.IMPPlanningUseCase, validate *validator.Validate, userHelper helper.IUserHelper) IMPPlanningHandler {
	return &MPPlanningHandler{
		Log:        log,
		Viper:      viper,
		UseCase:    useCase,
		Validate:   validate,
		UserHelper: userHelper,
	}
}

func MPPlanningHandlerFactory(log *logrus.Logger, viper *viper.Viper) IMPPlanningHandler {
	usecase := usecase.MPPlanningUseCaseFactory(viper, log)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	return NewMPPlanningHandler(log, viper, usecase, validate, userHelper)
}

func (h *MPPlanningHandler) FindAllHeadersPaginated(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := ctx.Query("search")
	if search == "" {
		search = ""
	}

	req := request.FindAllHeadersPaginatedMPPlanningRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	resp, err := h.UseCase.FindAllHeadersPaginated(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.FindAllHeadersPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find all headers paginated success", resp)
}

func (h *MPPlanningHandler) FindAllHeadersByRequestorIDPaginated(ctx *gin.Context) {
	user, err := middleware.GetUser(ctx, h.Log)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		h.Log.Errorf("Error when getting user: %v", err)
		return
	}
	if user == nil {
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		h.Log.Errorf("User not found")
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := ctx.Query("search")
	if search == "" {
		search = ""
	}

	req := request.FindAllHeadersPaginatedMPPlanningRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	userData, ok := user["user"].(map[string]interface{})
	if !ok {
		h.Log.Errorf("User information is missing or invalid")
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "User information is missing or invalid")
		return
	}

	userID, ok := userData["id"].(string)
	if !ok {
		h.Log.Errorf("[MPPlanningHandler.FindAllHeadersByRequestorIDPaginated] User ID is missing or invalid")
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "User ID is missing or invalid")
		return
	}

	requestorID, err := uuid.Parse(userID)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.FindAllHeadersByRequestorIDPaginated] Invalid User ID format: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "Invalid User ID format")
		return
	}

	resp, err := h.UseCase.FindAllHeadersByRequestorIDPaginated(requestorID, &req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.FindAllHeadersByRequestorIDPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find all headers by requestor id paginated success", resp)
}

func (h *MPPlanningHandler) GenerateDocumentNumber(ctx *gin.Context) {
	dateNow := time.Now()

	resp, err := h.UseCase.GenerateDocumentNumber(dateNow)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.GenerateDocumentNumber] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "generate document number success", resp)
}

func (h *MPPlanningHandler) FindById(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	req := request.FindHeaderByIdMPPlanningRequest{
		ID: id,
	}

	resp, err := h.UseCase.FindById(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.FindById] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find by id success", resp)
}

func (h *MPPlanningHandler) Create(ctx *gin.Context) {
	// Parse multipart form data
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
		h.Log.Errorf("[MPPlanningHandler.Create] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	// Get JSON payload from form data
	jsonData := ctx.Request.FormValue("payload")
	if jsonData == "" {
		h.Log.Errorf("[MPPlanningHandler.Create] JSON payload is empty")
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "JSON payload is empty")
		return
	}

	payload := new(request.CreateHeaderMPPlanningRequest)
	if err := json.Unmarshal([]byte(jsonData), payload); err != nil {
		h.Log.Errorf("[MPPlanningHandler.Create] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	// Validate payload
	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[MPPlanningHandler.Create] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	// Get user information
	user, err := middleware.GetUser(ctx, h.Log)
	if err != nil {
		h.Log.Errorf("Error when getting user: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}
	if user == nil {
		h.Log.Errorf("User not found")
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		return
	}

	// check if organization location exists in user
	organizationLocationID, err := h.UserHelper.CheckOrganizationLocation(user)
	if err != nil {
		h.Log.Errorf("Error when checking organization location: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	if organizationLocationID != payload.OrganizationLocationID {
		h.Log.Errorf("Organization location ID is not match")
		utils.ErrorResponse(ctx, 400, "error", "Organization location ID is not match")
		return
	}

	// Process uploaded files
	form, err := ctx.MultipartForm()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		h.Log.Errorf("[MPPlanningHandler.Create] " + err.Error())
		return
	}

	files := form.File["attachments"]
	var attachments []request.ManpowerAttachmentRequest

	for _, file := range files {
		// Generate a new file name with a timestamp
		timestamp := time.Now().UnixNano()
		extension := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", timestamp, extension)
		filePath := "storage/mp_planning_header/attachments/" + newFileName

		// Save the file or process it as needed
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			h.Log.Errorf("[MPPlanningHandler.Create] " + err.Error())
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
			return
		} else {
			h.Log.Infof("File %s uploaded successfully", filePath)
		}

		// Get the file type (MIME type)
		fileType := file.Header.Get("Content-Type")

		// Add file information to attachments
		attachments = append(attachments, request.ManpowerAttachmentRequest{
			FileName: newFileName,
			FilePath: filePath, // Or generate a URL if needed
			FileType: fileType,
		})
	}

	// Add attachments to payload
	payload.Attachments = attachments

	// Call use case to create the record
	resp, err := h.UseCase.Create(payload)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		h.Log.Errorf("[MPPlanningHandler.Create] " + err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "create success", resp)
}

func (h *MPPlanningHandler) UpdateStatusMPPPlanningHeader(ctx *gin.Context) {
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
		h.Log.Errorf("[MPPlanningHandler.UpdateStatusMPPPlanningHeader] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	// Get JSON payload from form data
	jsonData := ctx.Request.FormValue("payload")
	if jsonData == "" {
		h.Log.Errorf("[MPPlanningHandler.UpdateStatusMPPPlanningHeader] JSON payload is empty")
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "JSON payload is empty")
		return
	}

	payload := new(request.UpdateStatusMPPlanningHeaderRequest)
	if err := json.Unmarshal([]byte(jsonData), payload); err != nil {
		h.Log.Errorf("[MPPlanningHandler.UpdateStatusMPPPlanningHeader] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	} else {
		h.Log.Infof("Payload: %v", payload)
	}

	// Validate payload
	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[MPPlanningHandler.UpdateStatusMPPPlanningHeader] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	// Process uploaded files
	form, err := ctx.MultipartForm()
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.UpdateStatusMPPPlanningHeader] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	files := form.File["attachments"]
	var attachments []request.ManpowerAttachmentRequest

	for _, file := range files {
		// Generate a new file name with a timestamp
		timestamp := time.Now().UnixNano()
		extension := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", timestamp, extension)
		filePath := "storage/mp_planning_header/attachments/" + newFileName

		// Save the file or process it as needed
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			h.Log.Errorf("[MPPlanningHandler.UpdateStatusMPPPlanningHeader] " + err.Error())
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
			return
		} else {
			h.Log.Infof("File %s uploaded successfully", filePath)
		}

		// Get the file type (MIME type)
		fileType := file.Header.Get("Content-Type")

		// Add file information to attachments
		attachments = append(attachments, request.ManpowerAttachmentRequest{
			FileName: newFileName,
			FilePath: filePath, // Or generate a URL if needed
			FileType: fileType,
		})
	}

	// Add attachments to payload
	payload.Attachments = attachments

	err = h.UseCase.UpdateStatusMPPlanningHeader(payload)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.UpdateStatusMPPPlanningHeader] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "update status success", nil)
}

func (h *MPPlanningHandler) GetPlanningApprovalHistoryByHeaderId(ctx *gin.Context) {
	headerId := ctx.Param("header_id")

	if headerId == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "header_id is required")
		return
	}

	resp, err := h.UseCase.GetPlanningApprovalHistoryByHeaderId(uuid.MustParse(headerId))
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.GetPlanningApprovalHistoryByHeaderId] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "get planning approval history by header id success", resp)
}

func (h *MPPlanningHandler) Update(ctx *gin.Context) {
	// Parse multipart form data
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
		h.Log.Errorf("[MPPlanningHandler.Update file] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	// Get JSON payload from form data
	jsonData := ctx.Request.FormValue("payload")
	if jsonData == "" {
		h.Log.Errorf("[MPPlanningHandler.Update] JSON payload is empty")
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "JSON payload is empty")
		return
	}

	payload := new(request.UpdateHeaderMPPlanningRequest)
	if err := json.Unmarshal([]byte(jsonData), payload); err != nil {
		h.Log.Errorf("[MPPlanningHandler.Update] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	// Validate payload
	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[MPPlanningHandler.Update] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	// Get user information
	user, err := middleware.GetUser(ctx, h.Log)
	if err != nil {
		h.Log.Errorf("Error when getting user: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}
	if user == nil {
		h.Log.Errorf("User not found")
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		return
	}

	// check if organization location exists in user
	organizationLocationID, err := h.UserHelper.CheckOrganizationLocation(user)
	if err != nil {
		h.Log.Errorf("Error when checking organization location: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	if organizationLocationID != payload.OrganizationLocationID {
		h.Log.Errorf("Organization location ID is not match")
		utils.ErrorResponse(ctx, 400, "error", "Organization location ID is not match")
		return
	}

	// Process uploaded files
	form, err := ctx.MultipartForm()
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.Update] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	files := form.File["attachments"]
	var attachments []request.ManpowerAttachmentRequest

	for _, file := range files {
		// Generate a new file name with a timestamp
		timestamp := time.Now().UnixNano()
		extension := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", timestamp, extension)
		filePath := "storage/mp_planning_header/attachments/" + newFileName

		// Save the file or process it as needed
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			h.Log.Errorf("[MPPlanningHandler.Update] " + err.Error())
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
			return
		} else {
			h.Log.Infof("File %s uploaded successfully", filePath)
		}

		// Get the file type (MIME type)
		fileType := file.Header.Get("Content-Type")

		// Add file information to attachments
		attachments = append(attachments, request.ManpowerAttachmentRequest{
			FileName: newFileName,
			FilePath: filePath, // Or generate a URL if needed
			FileType: fileType,
		})
	}

	// Add attachments to payload
	payload.Attachments = attachments

	resp, err := h.UseCase.Update(payload)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.Update] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "update success", resp)
}

func (h *MPPlanningHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	req := request.DeleteHeaderMPPlanningRequest{
		ID: id,
	}

	err := h.UseCase.Delete(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.Delete] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "delete success", nil)
}

func (h *MPPlanningHandler) FindHeaderByMPPPeriodId(ctx *gin.Context) {
	mppPeriodId := ctx.Param("mpp_period_id")

	if mppPeriodId == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "mpp_period_id is required")
		return
	}

	req := request.FindHeaderByMPPPeriodIdMPPlanningRequest{
		MPPPeriodID: mppPeriodId,
	}

	resp, err := h.UseCase.FindHeaderByMPPPeriodId(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.FindHeaderByMPPPeriodId] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find header by mpp period id success", resp)
}

func (h *MPPlanningHandler) FindAllLinesByHeaderIdPaginated(ctx *gin.Context) {
	headerId := ctx.Param("header_id")

	if headerId == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "header_id is required")
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	req := request.FindAllLinesByHeaderIdPaginatedMPPlanningLineRequest{
		HeaderID: headerId,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := h.UseCase.FindAllLinesByHeaderIdPaginated(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.FindAllLinesByHeaderIdPaginated] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find all lines by header id paginated success", resp)
}

func (h *MPPlanningHandler) FindLineById(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	req := request.FindLineByIdMPPlanningLineRequest{
		ID: id,
	}

	resp, err := h.UseCase.FindLineById(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.FindLineById] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "find line by id success", resp.MPPlanningLine)
}

func (h *MPPlanningHandler) CreateLine(ctx *gin.Context) {
	var req request.CreateLineMPPlanningLineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[MPPlanningHandler.CreateLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[MPPlanningHandler.CreateLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	resp, err := h.UseCase.CreateLine(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.CreateLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "create line success", resp)
}

func (h *MPPlanningHandler) UpdateLine(ctx *gin.Context) {
	var req request.UpdateLineMPPlanningLineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[MPPlanningHandler.UpdateLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[MPPlanningHandler.UpdateLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	resp, err := h.UseCase.UpdateLine(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.UpdateLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "update line success", resp)
}

func (h *MPPlanningHandler) DeleteLine(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "id is required")
		return
	}

	req := request.DeleteLineMPPlanningLineRequest{
		ID: id,
	}

	err := h.UseCase.DeleteLine(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.DeleteLine] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "delete line success", nil)
}

func (h *MPPlanningHandler) CreateOrUpdateBatchLineMPPlanningLines(ctx *gin.Context) {
	var req request.CreateOrUpdateBatchLineMPPlanningLinesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[MPPlanningHandler.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[MPPlanningHandler.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	err := h.UseCase.CreateOrUpdateBatchLineMPPlanningLines(&req)
	if err != nil {
		h.Log.Errorf("[MPPlanningHandler.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "create or update batch line success", nil)
}
