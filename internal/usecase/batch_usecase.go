package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/dto"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IBatchUsecase interface {
	CreateBatchHeaderAndLines(req *request.CreateBatchHeaderAndLinesRequest) (*response.BatchResponse, error)
	FindByStatus(status entity.BatchHeaderApprovalStatus) (*response.BatchResponse, error)
	FindById(id string) (*response.BatchResponse, error)
	FindDocumentByID(id string) (*response.RealDocumentBatchResponse, error)
	FindByCurrentDocumentDateAndStatus(status entity.BatchHeaderApprovalStatus) (*response.BatchResponse, error)
}

type BatchUsecase struct {
	Viper      *viper.Viper
	Log        *logrus.Logger
	Repo       repository.IBatchRepository
	OrgMessage messaging.IOrganizationMessage
	EmpMessage messaging.IEmployeeMessage
	batchDTO   dto.IBatchDTO
}

func NewBatchUsecase(viper *viper.Viper, log *logrus.Logger, repo repository.IBatchRepository, orgMessage messaging.IOrganizationMessage, empMessage messaging.IEmployeeMessage, batchDTO dto.IBatchDTO) IBatchUsecase {
	return &BatchUsecase{
		Viper:      viper,
		Log:        log,
		Repo:       repo,
		OrgMessage: orgMessage,
		EmpMessage: empMessage,
		batchDTO:   batchDTO,
	}
}

func (uc *BatchUsecase) CreateBatchHeaderAndLines(req *request.CreateBatchHeaderAndLinesRequest) (*response.BatchResponse, error) {
	dateNow := time.Now()
	documentNumber := "MPP/BATCH/" + dateNow.Format("20060102") + "/001"

	foundBatchHeader, err := uc.Repo.GetHeadersByDocumentDate(dateNow.Format("2006-01-02"))
	if err != nil {
		uc.Log.Errorf("[MPPlanningUseCase.GenerateDocumentNumber] " + err.Error())
		return nil, err
	}

	if foundBatchHeader == nil {
		documentNumber = "MPP/BATCH/" + dateNow.Format("20060102") + "/001"
	} else {
		documentNumber = "MPP/BATCH/" + dateNow.Format("20060102") + "/" + fmt.Sprintf("%03d", len(*&foundBatchHeader)+1)
	}

	var batchHeader *entity.BatchHeader

	if req.DocumentNumber != "" {
		batchHeader = &entity.BatchHeader{
			DocumentNumber: req.DocumentNumber,
			DocumentDate:   dateNow,
			Status:         req.Status,
		}
	} else {
		batchHeader = &entity.BatchHeader{
			DocumentNumber: documentNumber,
			DocumentDate:   dateNow,
			Status:         entity.BatchHeaderApprovalStatusNeedApproval,
		}
	}

	batchLines := make([]entity.BatchLine, len(req.BatchLines))
	for i, bl := range req.BatchLines {
		if bl.OrganizationLocationID != "" {
			// Check if organization location exist
			orgLocExist, err := uc.OrgMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
				ID: bl.OrganizationLocationID,
			})

			if err != nil {
				uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
				return nil, err
			}

			if orgLocExist == nil {
				uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] Organization Location not found")
				return nil, errors.New("Organization Location not found")
			}
		}

		if bl.OrganizationID != "" {
			// Check if organization exist
			orgExist, err := uc.OrgMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
				ID: bl.OrganizationID,
			})

			if err != nil {
				uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] " + err.Error())
				return nil, err
			}

			if orgExist == nil {
				uc.Log.Errorf("[MPPlanningUseCase.CreateOrUpdateBatchLineMPPlanningLines] Organization not found")
				return nil, errors.New("Organization not found")
			}
		}

		batchLines[i] = entity.BatchLine{
			MPPlanningHeaderID:     uuid.MustParse(bl.MPPlanningHeaderID),
			OrganizationID:         func(u uuid.UUID) *uuid.UUID { return &u }(uuid.MustParse(bl.OrganizationID)),
			OrganizationLocationID: func(u uuid.UUID) *uuid.UUID { return &u }(uuid.MustParse(bl.OrganizationLocationID)),
		}
	}

	batchHeaderExists, err := uc.Repo.FindByStatus(entity.BatchHeaderApprovalStatusNeedApproval)
	if err != nil {
		uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] " + err.Error())
		return nil, err
	}

	if batchHeaderExists != nil {
		err = uc.Repo.InsertLinesByBatchHeaderID(batchHeaderExists.ID.String(), batchLines)
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] " + err.Error())
			return nil, err
		}

		err = uc.Repo.DeleteLinesNotInBatchLines(batchHeaderExists.ID.String(), batchLines)
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] " + err.Error())
			return nil, err
		}

		findBatchHeader, err := uc.Repo.FindById(batchHeaderExists.ID.String())
		if err != nil {
			uc.Log.Errorf("[BatchUsecase.CreateBatchHeaderAndLines] " + err.Error())
			return nil, err
		}

		return uc.batchDTO.ConvertBatchHeaderEntityToResponse(findBatchHeader), nil
	}

	resp, err := uc.Repo.CreateBatchHeaderAndLines(batchHeader, batchLines)
	if err != nil {
		return nil, err
	}

	return uc.batchDTO.ConvertBatchHeaderEntityToResponse(resp), nil
}

func (uc *BatchUsecase) FindByStatus(status entity.BatchHeaderApprovalStatus) (*response.BatchResponse, error) {
	resp, err := uc.Repo.FindByStatus(status)
	if err != nil {
		return nil, err
	}

	return uc.batchDTO.ConvertBatchHeaderEntityToResponse(resp), nil
}

func (uc *BatchUsecase) FindById(id string) (*response.BatchResponse, error) {
	resp, err := uc.Repo.FindById(id)
	if err != nil {
		return nil, err
	}

	return uc.batchDTO.ConvertBatchHeaderEntityToResponse(resp), nil
}

func (uc *BatchUsecase) FindDocumentByID(id string) (*response.RealDocumentBatchResponse, error) {
	resp, err := uc.Repo.FindById(id)
	if err != nil {
		return nil, err
	}

	return uc.batchDTO.ConvertRealDocumentBatchResponse(resp), nil
}

func (uc *BatchUsecase) FindByCurrentDocumentDateAndStatus(status entity.BatchHeaderApprovalStatus) (*response.BatchResponse, error) {
	resp, err := uc.Repo.FindByCurrentDocumentDateAndStatus(status)
	if err != nil {
		return nil, err
	}

	return uc.batchDTO.ConvertBatchHeaderEntityToResponse(resp), nil
}

func BatchUsecaseFactory(viper *viper.Viper, log *logrus.Logger) IBatchUsecase {
	repo := repository.BatchRepositoryFactory(log)
	orgMessage := messaging.OrganizationMessageFactory(log)
	empMessage := messaging.EmployeeMessageFactory(log)
	batchDTO := dto.BatchDTOFactory(log)
	return NewBatchUsecase(viper, log, repo, orgMessage, empMessage, batchDTO)
}
