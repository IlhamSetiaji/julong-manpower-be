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
	FindByJobId(request *request.FindByJobIdJobPlafonRequest) (*response.FindByJobIdJobPlafonResponse, error)
	Create(request *request.CreateJobPlafonRequest) (*response.CreateJobPlafonResponse, error)
	Update(request *request.UpdateJobPlafonRequest) (*response.UpdateJobPlafonResponse, error)
	Delete(request *request.DeleteJobPlafonRequest) error
}

type JobPlafonUseCase struct {
	Log                 *logrus.Logger
	JobPlafonRepository repository.IJobPlafonRepository
	JobPlafonMessage    messaging.IJobPlafonMessage
	JobMessage          messaging.IJobMessage
}

func NewJobPlafonUseCase(log *logrus.Logger, repo repository.IJobPlafonRepository, message messaging.IJobPlafonMessage, jm messaging.IJobMessage) IJobPlafonUseCase {
	return &JobPlafonUseCase{
		Log:                 log,
		JobPlafonRepository: repo,
		JobPlafonMessage:    message,
		JobMessage:          jm,
	}
}

func (uc *JobPlafonUseCase) FindAllPaginated(req *request.FindAllPaginatedJobPlafonRequest) (*response.FindAllPaginatedJobPlafonResponse, error) {
	jobPlafons, total, err := uc.JobPlafonRepository.FindAllPaginated(req.Page, req.PageSize, req.Search)
	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.FindAllPaginated] " + err.Error())
		return nil, err
	}

	for i, jobPlafon := range *jobPlafons {
		// messageResponse, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
		// 	ID: jobPlafon.JobID.String(),
		// })

		// if err != nil {
		// 	uc.Log.Errorf("[JobPlafonUseCase.FindAllPaginated Message] " + err.Error())
		// 	return nil, err
		// }

		// jobPlafon.JobName = messageResponse.Name
		messageResponse, err := uc.JobMessage.SendFindJobDataByIdMessage(request.SendFindJobByIDMessageRequest{
			ID: jobPlafon.JobID.String(),
		})

		if err != nil {
			uc.Log.Errorf("[JobPlafonUseCase.FindAllPaginated Message] " + err.Error())
			return nil, err
		}

		jobPlafon.JobName = messageResponse.Name
		jobPlafon.OrganizationName = messageResponse.OrganizationName

		(*jobPlafons)[i] = jobPlafon
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

func (uc *JobPlafonUseCase) FindByJobId(payload *request.FindByJobIdJobPlafonRequest) (*response.FindByJobIdJobPlafonResponse, error) {
	messageResponse, err := uc.JobPlafonMessage.SendCheckJobExistMessage(request.CheckJobExistMessageRequest{
		ID: payload.JobID,
	})

	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.FindByJobId Message] " + err.Error())
		return nil, err
	}

	if !messageResponse.Exist {
		uc.Log.Errorf("[JobPlafonUseCase.FindByJobId] Job not found")
		return nil, errors.New("job not found")
	}

	jobPlafon, err := uc.JobPlafonRepository.FindByJobId(uuid.MustParse(payload.JobID))
	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.FindByJobId] " + err.Error())
		return nil, err
	}

	jobResponse, err := uc.JobMessage.SendFindJobDataByIdMessage(request.SendFindJobByIDMessageRequest{
		ID: jobPlafon.JobID.String(),
	})

	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.FindAllPaginated Message] " + err.Error())
		return nil, err
	}

	jobPlafon.JobName = jobResponse.Name
	jobPlafon.OrganizationName = jobResponse.OrganizationName

	return &response.FindByJobIdJobPlafonResponse{
		JobPlafon: jobPlafon,
	}, nil
}

func (uc *JobPlafonUseCase) Create(payload *request.CreateJobPlafonRequest) (*response.CreateJobPlafonResponse, error) {
	jobExist, err := uc.JobPlafonRepository.FindByJobId(uuid.MustParse(payload.JobID))

	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.Create] " + err.Error())
		return nil, err
	}

	if jobExist != nil {
		uc.Log.Errorf("[JobPlafonUseCase.Create] Job plafon already exist")
		return nil, errors.New("job plafon already exist")
	}

	messageResponse, err := uc.JobPlafonMessage.SendCheckJobExistMessage(request.CheckJobExistMessageRequest{
		ID: payload.JobID,
	})

	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.Create Message] " + err.Error())
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
	jm := messaging.JobMessageFactory(log)
	return NewJobPlafonUseCase(log, repo, message, jm)
}
