package usecase

import (
	"errors"
	"strings"

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
	SyncJobPlafon() error
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

func (uc *JobPlafonUseCase) SyncJobPlafon() error {
	// get jobs data using rabbitmq
	jobPlafonMessageResponse, err := uc.JobMessage.SendGetAllJobDataMessage()
	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.SyncJobPlafon] " + err.Error())
		return err
	}

	// loop jobs data to find job plafon (create when doesnt exist)
	for _, job := range *jobPlafonMessageResponse {
		jobPlafon, err := uc.JobPlafonRepository.FindByJobId(job.ID)
		if err != nil {
			uc.Log.Errorf("[JobPlafonUseCase.SyncJobPlafon] " + err.Error())
			return err
		}

		if jobPlafon == nil {
			jobPlafonEntity := entity.JobPlafon{
				JobID:  &job.ID,
				Plafon: 0,
			}

			_, err := uc.JobPlafonRepository.Create(&jobPlafonEntity)
			if err != nil {
				uc.Log.Errorf("[JobPlafonUseCase.SyncJobPlafon] " + err.Error())
				return err
			}
		}
	}

	return nil
}

func (uc *JobPlafonUseCase) FindAllPaginated(req *request.FindAllPaginatedJobPlafonRequest) (*response.FindAllPaginatedJobPlafonResponse, error) {
	// get jobs ids using rabbitmq
	jobPlafonMessageResponse, err := uc.JobMessage.SendFindAllJobsByOrganizationIDMessage(req.OrganizationID)
	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.FindAllPaginated] " + err.Error())
		return nil, err
	}

	var jobIDs []string
	for _, job := range *jobPlafonMessageResponse {
		jobIDs = append(jobIDs, job.ID.String())
	}

	var filter map[string]interface{}
	if len(jobIDs) > 0 {
		filter = map[string]interface{}{
			"job_ids": jobIDs,
		}
	}

	jobPlafons, total, err := uc.JobPlafonRepository.FindAllPaginated(req.Page, req.PageSize, req.Search, filter)
	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.FindAllPaginated] " + err.Error())
		return nil, err
	}

	var filteredJobPlafons []entity.JobPlafon

	for _, jobPlafon := range *jobPlafons {
		messageResponse, err := uc.JobMessage.SendFindJobDataByIdMessage(request.SendFindJobByIDMessageRequest{
			ID: jobPlafon.JobID.String(),
		})

		if err != nil {
			uc.Log.Errorf("[JobPlafonUseCase.FindAllPaginated Message] " + err.Error())
			return nil, err
		}

		jobPlafon.JobName = messageResponse.Name
		jobPlafon.OrganizationName = messageResponse.OrganizationName

		if req.Search != "" {
			// if jobPlafon.OrganizationName contains req.Search (case-insensitive)
			if !strings.Contains(strings.ToLower(jobPlafon.OrganizationName), strings.ToLower(req.Search)) && !strings.Contains(strings.ToLower(jobPlafon.JobName), strings.ToLower(req.Search)) {
				continue
			}
		}

		filteredJobPlafons = append(filteredJobPlafons, jobPlafon)
	}

	// Update the original slice with the filtered results
	*jobPlafons = filteredJobPlafons

	return &response.FindAllPaginatedJobPlafonResponse{
		JobPlafons: jobPlafons,
		Total:      total,
	}, nil
}

func (uc *JobPlafonUseCase) FindById(req *request.FindByIdJobPlafonRequest) (*response.FindByIdJobPlafonResponse, error) {
	jobPlafon, err := uc.JobPlafonRepository.FindById(uuid.MustParse(req.ID))
	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.FindById] " + err.Error())
		return nil, err
	}

	jobResponse, err := uc.JobMessage.SendFindJobDataByIdMessage(request.SendFindJobByIDMessageRequest{
		ID: jobPlafon.JobID.String(),
	})

	if err != nil {
		uc.Log.Errorf("[JobPlafonUseCase.FindAllPaginated Message] " + err.Error())
		return nil, err
	}

	uc.Log.Infof("jobResponse: %+v", jobResponse.Name)

	jobPlafon.JobName = jobResponse.Name
	jobPlafon.OrganizationName = jobResponse.OrganizationName

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

	uc.Log.Infof("jobResponse: %+v", jobResponse.Name)

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
