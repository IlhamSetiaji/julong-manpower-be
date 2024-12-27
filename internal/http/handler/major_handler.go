package handler

import (
	"net/http"

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

type IMajorHandler interface {
	FindAll(ctx *gin.Context)
	FindById(ctx *gin.Context)
	GetMajorsByEducationLevel(ctx *gin.Context)
}

type MajorHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Uc       usecase.IMajorUsecase
	Validate *validator.Validate
}

func NewMajorHandler(log *logrus.Logger, viper *viper.Viper, uc usecase.IMajorUsecase, validate *validator.Validate) *MajorHandler {
	return &MajorHandler{Log: log, Viper: viper, Uc: uc, Validate: validate}
}

func (h *MajorHandler) FindAll(ctx *gin.Context) {
	majors, err := h.Uc.FindAll()
	if err != nil {
		h.Log.Errorf("[MajorHandler.FindAll] error when finding all majors. Error: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when finding all majors", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success finding all majors", majors)
}

func (h *MajorHandler) FindById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "ID is required", "ID is required")
		return
	}

	request := request.FindByIdMajorRequest{ID: uuid.MustParse(id)}
	err := h.Validate.Struct(request)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	major, err := h.Uc.FindById(request)
	if err != nil {
		h.Log.Errorf("[MajorHandler.FindById] error when finding major by id. Error: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when finding major by id", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success finding major by id", major)
}

func (h *MajorHandler) GetMajorsByEducationLevel(ctx *gin.Context) {
	educationLevel := ctx.Param("education_level")
	if educationLevel == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Education level is required", "Education level is required")
		return
	}

	majors, err := h.Uc.GetMajorsByEducationLevel(educationLevel)
	if err != nil {
		h.Log.Errorf("[MajorHandler.GetMajorsByEducationLevel] error when getting majors by education level. Error: %v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Error when getting majors by education level", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Success getting majors by education level", majors)
}

func MajorHandlerFactory(log *logrus.Logger, viper *viper.Viper) *MajorHandler {
	uc := usecase.MajorUsecaseFactory(log)
	validate := config.NewValidator(viper)
	return NewMajorHandler(log, viper, uc, validate)
}
