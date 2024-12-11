package usecase

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/dto"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IMPRequestUseCase interface {
	Create(req *request.CreateMPRequestHeaderRequest) (*response.MPRequestHeaderResponse, error)
}

type MPRequestUseCase struct {
	Viper                  *viper.Viper
	Log                    *logrus.Logger
	MPRequestRepository    repository.IMPRequestRepository
	RequestMajorRepository repository.IRequestMajorRepository
	OrganizationMessage    messaging.IOrganizationMessage
	JobPlafonMessage       messaging.IJobPlafonMessage
	UserMessage            messaging.IUserMessage
}

func NewMPRequestUseCase(
	viper *viper.Viper,
	log *logrus.Logger,
	mpRequestRepository repository.IMPRequestRepository,
	requestMajorRepository repository.IRequestMajorRepository,
	organizationMessage messaging.IOrganizationMessage,
	jobPlafonMessage messaging.IJobPlafonMessage,
	userMessage messaging.IUserMessage,
) IMPRequestUseCase {
	return &MPRequestUseCase{
		Viper:                  viper,
		Log:                    log,
		MPRequestRepository:    mpRequestRepository,
		RequestMajorRepository: requestMajorRepository,
		OrganizationMessage:    organizationMessage,
		JobPlafonMessage:       jobPlafonMessage,
		UserMessage:            userMessage,
	}
}

func (uc *MPRequestUseCase) Create(req *request.CreateMPRequestHeaderRequest) (*response.MPRequestHeaderResponse, error) {
	// check if organization is exist
	orgExist, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: req.OrganizationID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] error when send find organization by id message: %v", err)
		return nil, err
	}

	if orgExist == nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] organization with id %s is not exist", req.OrganizationID.String())
		return nil, errors.New("organization is not exist")
	}

	// check if organization location is exist
	orgLocExist, err := uc.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
		ID: req.OrganizationLocationID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] error when send find organization location by id message: %v", err)
		return nil, err
	}

	if orgLocExist == nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] organization location with id %s is not exist", req.OrganizationLocationID.String())
		return nil, errors.New("organization location is not exist")
	}

	// check if for organization is exist
	forOrgExist, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: req.ForOrganizationID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] error when send find for organization by id message: %v", err)
		return nil, err
	}

	if forOrgExist == nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] for organization with id %s is not exist", req.ForOrganizationID.String())
		return nil, errors.New("for organization is not exist")
	}

	// check if for organization location is exist
	forOrgLocExist, err := uc.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
		ID: req.ForOrganizationLocationID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] error when send find for organization location by id message: %v", err)
		return nil, err
	}

	if forOrgLocExist == nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] for organization location with id %s is not exist", req.ForOrganizationLocationID.String())
		return nil, errors.New("for organization location is not exist")
	}

	// check if for organization structure is exist
	forOrgStructExist, err := uc.OrganizationMessage.SendFindOrganizationStructureByIDMessage(request.SendFindOrganizationStructureByIDMessageRequest{
		ID: req.ForOrganizationStructureID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] error when send find for organization structure by id message: %v", err)
		return nil, err
	}

	if forOrgStructExist == nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] for organization structure with id %s is not exist", req.ForOrganizationStructureID.String())
		return nil, errors.New("for organization structure is not exist")
	}

	// check if job ID is exist
	jobExist, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
		ID: req.JobID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] error when send find job by id message: %v", err)
		return nil, err
	}

	if jobExist == nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] job with id %s is not exist", req.JobID.String())
		return nil, errors.New("job is not exist")
	}

	// check if requestor ID is exist
	requestorExist, err := uc.UserMessage.SendFindUserByIDMessage(request.SendFindUserByIDMessageRequest{
		ID: req.RequestorID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] error when send find user by id message: %v", err)
		return nil, err
	}

	if requestorExist == nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] requestor with id %s is not exist", req.RequestorID.String())
		return nil, errors.New("requestor is not exist")
	}

	// check if department head is exist
	deptHeadExist, err := uc.UserMessage.SendFindUserByIDMessage(request.SendFindUserByIDMessageRequest{
		ID: req.DepartmentHead.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] error when send find user by id message: %v", err)
		return nil, err
	}

	if deptHeadExist == nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] department head with id %s is not exist", req.DepartmentHead.String())
		return nil, errors.New("department head is not exist")
	}

	mpRequestHeader, err := uc.MPRequestRepository.Create(dto.ConvertToEntity(req))
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] error when create mp request header: %v", err)
		return nil, err
	}

	// create request major
	for _, majorID := range req.MajorIDs {
		reqMajor, err := uc.RequestMajorRepository.Create(&entity.RequestMajor{
			MPRequestHeaderID: mpRequestHeader.ID,
			MajorID:           majorID,
		})
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] error when create request major: %v", err)
			return nil, err
		}

		mpRequestHeader.RequestMajors = append(mpRequestHeader.RequestMajors, *reqMajor)
	}

	mpRequestHeader.OrganizationName = orgExist.Name
	mpRequestHeader.OrganizationLocationName = orgLocExist.Name
	mpRequestHeader.ForOrganizationName = forOrgExist.Name
	mpRequestHeader.ForOrganizationLocation = forOrgLocExist.Name
	mpRequestHeader.ForOrganizationStructure = forOrgStructExist.Name
	mpRequestHeader.JobName = jobExist.Name
	mpRequestHeader.RequestorName = requestorExist.Name
	mpRequestHeader.DepartmentHeadName = deptHeadExist.Name

	return dto.ConvertToResponse(mpRequestHeader), nil
}

func MPRequestUseCaseFactory(viper *viper.Viper, log *logrus.Logger) IMPRequestUseCase {
	mpRequestRepository := repository.MPRequestRepositoryFactory(log)
	requestMajorRepository := repository.RequestMajorRepositoryFactory(log)
	organizationMessage := messaging.OrganizationMessageFactory(log)
	jobPlafonMessage := messaging.JobPlafonMessageFactory(log)
	userMessage := messaging.UserMessageFactory(log)
	return NewMPRequestUseCase(viper, log, mpRequestRepository, requestMajorRepository, organizationMessage, jobPlafonMessage, userMessage)
}
