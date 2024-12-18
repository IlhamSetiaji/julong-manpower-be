package usecase

import (
	"errors"

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
	batchHeader := &entity.BatchHeader{
		DocumentNumber: req.DocumentNumber,
		Status:         req.Status,
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

func BatchUsecaseFactory(viper *viper.Viper, log *logrus.Logger) IBatchUsecase {
	repo := repository.BatchRepositoryFactory(log)
	orgMessage := messaging.OrganizationMessageFactory(log)
	empMessage := messaging.EmployeeMessageFactory(log)
	batchDTO := dto.BatchDTOFactory(log)
	return NewBatchUsecase(viper, log, repo, orgMessage, empMessage, batchDTO)
}
