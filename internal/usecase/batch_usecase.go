package usecase

import (
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
}

type BatchUsecase struct {
	Viper      *viper.Viper
	Log        *logrus.Logger
	Repo       repository.IBatchRepository
	OrgMessage messaging.IOrganizationMessage
	EmpMessage messaging.IEmployeeMessage
}

func NewBatchUsecase(viper *viper.Viper, log *logrus.Logger, repo repository.IBatchRepository, orgMessage messaging.IOrganizationMessage, empMessage messaging.IEmployeeMessage) IBatchUsecase {
	return &BatchUsecase{
		Viper:      viper,
		Log:        log,
		Repo:       repo,
		OrgMessage: orgMessage,
		EmpMessage: empMessage,
	}
}

func (uc *BatchUsecase) CreateBatchHeaderAndLines(req *request.CreateBatchHeaderAndLinesRequest) (*response.BatchResponse, error) {
	batchHeader := &entity.BatchHeader{
		DocumentNumber: req.DocumentNumber,
		Status:         req.Status,
	}

	batchLines := make([]entity.BatchLine, len(req.BatchLines))
	for i, bl := range req.BatchLines {
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

	return dto.ConvertBatchHeaderEntityToResponse(resp), nil
}

func BatchUsecaseFactory(viper *viper.Viper, log *logrus.Logger) IBatchUsecase {
	repo := repository.BatchRepositoryFactory(log)
	orgMessage := messaging.OrganizationMessageFactory(log)
	empMessage := messaging.EmployeeMessageFactory(log)
	return NewBatchUsecase(viper, log, repo, orgMessage, empMessage)
}
