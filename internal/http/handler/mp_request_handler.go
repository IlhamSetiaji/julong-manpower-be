package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/config"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/usecase"
	"github.com/IlhamSetiaji/julong-manpower-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IMPRequestHandler interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	UpdateStatusMPRequestHeader(ctx *gin.Context)
}

type MPRequestHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	UseCase  usecase.IMPRequestUseCase
	Validate *validator.Validate
}

func NewMPRequestHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	useCase usecase.IMPRequestUseCase,
	validate *validator.Validate,
) IMPRequestHandler {
	return &MPRequestHandler{
		Log:      log,
		Viper:    viper,
		UseCase:  useCase,
		Validate: validate,
	}
}

func MPRequestHandlerFactory(log *logrus.Logger, viper *viper.Viper) IMPRequestHandler {
	useCase := usecase.MPRequestUseCaseFactory(viper, log)
	validate := config.NewValidator(viper)
	return NewMPRequestHandler(log, viper, useCase, validate)
}

func (h *MPRequestHandler) FindAllPaginated(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := ctx.Query("search")
	if search == "" {
		search = ""
	}

	filter := make(map[string]interface{})

	status := ctx.Query("status")
	if status != "" {
		filter["status"] = status
	}

	departmentHead := ctx.Query("department_head")
	if departmentHead != "" {
		filter["department_head"] = departmentHead
	}

	vpGmDirector := ctx.Query("vp_gm_director")
	if vpGmDirector != "" {
		filter["vp_gm_director"] = vpGmDirector
	}

	ceo := ctx.Query("ceo")
	if ceo != "" {
		filter["ceo"] = ceo
	}

	hrdHoUnit := ctx.Query("hrd_ho_unit")
	if hrdHoUnit != "" {
		filter["hrd_ho_unit"] = hrdHoUnit
	}

	res, err := h.UseCase.FindAllPaginated(page, pageSize, search, filter)
	if err != nil {
		h.Log.Errorf("[MPRequestHandler.FindAllPaginated] error when find all paginated: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all paginated", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "MP Request Headers found", res)
}

func (h *MPRequestHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.Log.Errorf("[MPRequestHandler.FindByID] error when get ID from request")
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request", "ID is required")
		return
	}

	res, err := h.UseCase.FindByID(uuid.MustParse(id))
	if err != nil {
		h.Log.Errorf("[MPRequestHandler.FindByID] error when find by ID: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find by ID", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "MP Request Header found", res)
}

func (h *MPRequestHandler) Create(ctx *gin.Context) {
	var req request.CreateMPRequestHeaderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[MPRequestHandler.Create] error when bind json: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[MPRequestHandler.Create] error when validate request: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	res, err := h.UseCase.Create(&req)
	if err != nil {
		h.Log.Errorf("[MPRequestHandler.Create] error when create mp request header: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create mp request header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "MP Request Header created", res)
}

func (h *MPRequestHandler) Update(ctx *gin.Context) {
	var req request.CreateMPRequestHeaderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Errorf("[MPRequestHandler.Update] error when bind json: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if req.ID == "" {
		h.Log.Errorf("[MPRequestHandler.Update] error when get ID from request")
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", "ID is required")
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Errorf("[MPRequestHandler.Update] error when validate request: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	res, err := h.UseCase.Update(&req)
	if err != nil {
		h.Log.Errorf("[MPRequestHandler.Update] error when update mp request header: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update mp request header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "MP Request Header updated", res)
}

func (h *MPRequestHandler) UpdateStatusMPRequestHeader(ctx *gin.Context) {
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max memory
		h.Log.Errorf("[MPRequestHandler.UpdateStatusMPRequestHeader] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	}

	// Get JSON payload
	jsonData := ctx.Request.FormValue("payload")
	if jsonData == "" {
		h.Log.Errorf("[MPRequestHandler.UpdateStatusMPRequestHeader] error when get json payload")
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", "Invalid request body")
		return
	}

	payload := new(request.UpdateMPRequestHeaderRequest)
	if err := json.Unmarshal([]byte(jsonData), payload); err != nil {
		h.Log.Errorf("[MPRequestHandler.UpdateStatusMPRequestHeader] error when unmarshal json payload: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "error", err.Error())
		return
	} else {
		h.Log.Infof("[MPRequestHandler.UpdateStatusMPRequestHeader] payload: %v", payload)
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("[MPRequestHandler.UpdateStatusMPRequestHeader] error when validate request: %v", err)
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// process attachments
	form, err := ctx.MultipartForm()
	if err != nil {
		h.Log.Errorf("[MPRequestHandler.UpdateStatusMPRequestHeader] error when get multipart form: %v", err)
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
		filePath := "storage/mp_request_header/attachments/" + newFileName

		// save the file
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			h.Log.Errorf("[MPRequestHandler.UpdateStatusMPRequestHeader] error when save uploaded file: %v", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
			return
		} else {
			h.Log.Infof("[MPRequestHandler.UpdateStatusMPRequestHeader] file saved to: %s", filePath)
		}

		// get the file type
		fileType := file.Header.Get("Content-Type")

		// add file information to the attachments
		attachments = append(attachments, request.ManpowerAttachmentRequest{
			FilePath: filePath,
			FileType: fileType,
			FileName: newFileName,
		})
	}

	payload.Attachments = attachments

	err = h.UseCase.UpdateStatusHeader(payload)
	if err != nil {
		h.Log.Errorf("[MPRequestHandler.UpdateStatusMPRequestHeader] error when update status: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update status", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "MP Request Header status updated", nil)
}
