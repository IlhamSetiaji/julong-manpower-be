package usecase

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IJobPlafonUseCase interface {
	FindAllPaginated(request *request.FindAllPaginatedJobPlafonRequest) (*response.FindAllPaginatedJobPlafonResponse, error)
	FindById(request *request.FindByIdJobPlafonRequest) (*response.FindByIdJobPlafonResponse, error)
	Create(request *request.CreateJobPlafonRequest) (*response.CreateJobPlafonResponse, error)
	Update(request *request.UpdateJobPlafonRequest) (*response.UpdateJobPlafonResponse, error)
	Delete(request *request.DeleteJobPlafonRequest) error
}

type JobPlafonUseCase struct {
	Log                 *logrus.Logger
	JobPlafonRepository repository.IJobPlafonRepository
	JobPlafonMessage    messaging.IJobPlafonMessage
}

func NewJobPlafonUseCase(log *logrus.Logger, repo repository.IJobPlafonRepository, message messaging.IJobPlafonMessage) IJobPlafonUseCase {
	return &JobPlafonUseCase{
		Log:                 log,
		JobPlafonRepository: repo,
		JobPlafonMessage:    message,
	}
}

func (uc *JobPlafonUseCase) FindAllPaginated(request *request.FindAllPaginatedJobPlafonRequest) (*response.FindAllPaginatedJobPlafonResponse, error) {
	jobPlafons, total, err := uc.JobPlafonRepository.FindAllPaginated(request.Page, request.PageSize, request.Search)
	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.FindAllPaginated] " + err.Error())
		return nil, err
	}

	return &response.FindAllPaginatedJobPlafonResponse{
		JobPlafons: jobPlafons,
		Total:      total,
	}, nil
}

func (uc *JobPlafonUseCase) FindById(request *request.FindByIdJobPlafonRequest) (*response.FindByIdJobPlafonResponse, error) {
	jobPlafon, err := uc.JobPlafonRepository.FindById(uuid.MustParse(request.ID))
	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.FindById] " + err.Error())
		return nil, err
	}

	return &response.FindByIdJobPlafonResponse{
		JobPlafon: jobPlafon,
	}, nil
}

func (uc *JobPlafonUseCase) Create(payload *request.CreateJobPlafonRequest) (*response.CreateJobPlafonResponse, error) {
	messageResponse, err := uc.JobPlafonMessage.SendCheckJobExistMessage(request.CheckJobExistMessageRequest{
		ID: payload.JobID,
	})

	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.Create] " + err.Error())
		return nil, err
	}

	if !messageResponse.Exist {
		uc.Log.Errorf("[JobPlafonUseCase.Create] Job not found")
		return nil, errors.New("job not found")
	}

	jobPlafonEntity := entity.JobPlafon{
		JobID:  func(id string) *uuid.UUID { u := uuid.MustParse(id); return &u }(payload.JobID),
		Plafon: payload.Plafon,
	}

	jobPlafon, err := uc.JobPlafonRepository.Create(&jobPlafonEntity)
	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.Create] " + err.Error())
		return nil, err
	}

	return &response.CreateJobPlafonResponse{
		JobPlafon: jobPlafon,
	}, nil
}

func (uc *JobPlafonUseCase) Update(payload *request.UpdateJobPlafonRequest) (*response.UpdateJobPlafonResponse, error) {
	messageResponse, err := uc.JobPlafonMessage.SendCheckJobExistMessage(request.CheckJobExistMessageRequest{
		ID: payload.JobID,
	})

	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.Create] " + err.Error())
		return nil, err
	}

	if !messageResponse.Exist {
		uc.Log.Errorf("[JobPlafonUseCase.Create] Job not found")
		return nil, errors.New("job not found")
	}

	jobPlafonEntity := &entity.JobPlafon{
		ID:     uuid.MustParse(payload.ID),
		JobID:  func(id string) *uuid.UUID { u := uuid.MustParse(id); return &u }(payload.JobID),
		Plafon: payload.Plafon,
	}
	jobPlafon, err := uc.JobPlafonRepository.Update(jobPlafonEntity)
	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.Update] " + err.Error())
		return nil, err
	}

	return &response.UpdateJobPlafonResponse{
		JobPlafon: jobPlafon,
	}, nil
}

func (uc *JobPlafonUseCase) Delete(request *request.DeleteJobPlafonRequest) error {
	err := uc.JobPlafonRepository.Delete(uuid.MustParse(request.ID))
	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.Delete] " + err.Error())
		return err
	}

	return nil
}

func JobPlafonUseCaseFactory(log *logrus.Logger) IJobPlafonUseCase {
	repo := repository.JobPlafonRepositoryFactory(log)
	message := messaging.JobPlafonMessageFactory(log)
	return NewJobPlafonUseCase(log, repo, message)
}
