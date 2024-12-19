package usecase

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/dto"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/helper"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IMPRequestUseCase interface {
	Create(req *request.CreateMPRequestHeaderRequest) (*response.MPRequestHeaderResponse, error)
	Update(req *request.CreateMPRequestHeaderRequest) (*response.MPRequestHeaderResponse, error)
	FindAllPaginated(page int, pageSize int, search string, filter map[string]interface{}) (*response.MPRequestPaginatedResponse, error)
	UpdateStatusHeader(req *request.UpdateMPRequestHeaderRequest) error
}

type MPRequestUseCase struct {
	Viper                  *viper.Viper
	Log                    *logrus.Logger
	MPRequestRepository    repository.IMPRequestRepository
	RequestMajorRepository repository.IRequestMajorRepository
	OrganizationMessage    messaging.IOrganizationMessage
	JobPlafonMessage       messaging.IJobPlafonMessage
	UserMessage            messaging.IUserMessage
	MPPPeriodRepo          repository.IMPPPeriodRepository
	EmpMessage             messaging.IEmployeeMessage
	MPRequestHelper        helper.IMPRequestHelper
}

func NewMPRequestUseCase(
	viper *viper.Viper,
	log *logrus.Logger,
	mpRequestRepository repository.IMPRequestRepository,
	requestMajorRepository repository.IRequestMajorRepository,
	organizationMessage messaging.IOrganizationMessage,
	jobPlafonMessage messaging.IJobPlafonMessage,
	userMessage messaging.IUserMessage,
	mppPeriodRepo repository.IMPPPeriodRepository,
	em messaging.IEmployeeMessage,
	mprHelper helper.IMPRequestHelper,
) IMPRequestUseCase {
	return &MPRequestUseCase{
		Viper:                  viper,
		Log:                    log,
		MPRequestRepository:    mpRequestRepository,
		RequestMajorRepository: requestMajorRepository,
		OrganizationMessage:    organizationMessage,
		JobPlafonMessage:       jobPlafonMessage,
		UserMessage:            userMessage,
		MPPPeriodRepo:          mppPeriodRepo,
		EmpMessage:             em,
		MPRequestHelper:        mprHelper,
	}
}

func (uc *MPRequestUseCase) Create(req *request.CreateMPRequestHeaderRequest) (*response.MPRequestHeaderResponse, error) {
	// check if mpp period exist
	mppPeriod, err := uc.MPPPeriodRepo.FindById(*req.MPPPeriodID)
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] error when find mpp period by id: %v", err)
		return nil, err
	}

	if mppPeriod == nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] mpp period with id %s is not exist", req.MPPPeriodID.String())
		return nil, errors.New("mpp period is not exist")
	}

	// check portal data
	portalResponse, err := uc.MPRequestHelper.CheckPortalData(req)
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] error when check portal data: %v", err)
		return nil, err
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

	mpRequestHeader.OrganizationName = portalResponse.OrganizationName
	mpRequestHeader.OrganizationLocationName = portalResponse.OrganizationLocationName
	mpRequestHeader.ForOrganizationName = portalResponse.ForOrganizationName
	mpRequestHeader.ForOrganizationLocation = portalResponse.ForOrganizationLocationName
	mpRequestHeader.ForOrganizationStructure = portalResponse.ForOrganizationStructureName
	mpRequestHeader.JobName = portalResponse.JobName
	mpRequestHeader.RequestorName = portalResponse.RequestorName
	mpRequestHeader.DepartmentHeadName = portalResponse.DepartmentHeadName
	mpRequestHeader.EmpOrganizationName = portalResponse.EmpOrganizationName
	mpRequestHeader.JobLevelName = portalResponse.JobLevelName
	mpRequestHeader.JobLevel = portalResponse.JobLevel
	mpRequestHeader.MPPPeriod = *mppPeriod

	return dto.ConvertToResponse(mpRequestHeader), nil
}

func (uc *MPRequestUseCase) Update(req *request.CreateMPRequestHeaderRequest) (*response.MPRequestHeaderResponse, error) {
	// check if mp request header is exist
	mpRequestHeaderExist, err := uc.MPRequestRepository.FindById(uuid.MustParse(req.ID))
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Update] error when find mp request header by id: %v", err)
		return nil, err
	}

	if mpRequestHeaderExist == nil {
		uc.Log.Errorf("[MPRequestUseCase.Update] mp request header with id %s is not exist", req.ID)
		return nil, errors.New("mp request header is not exist")
	}

	// check if mpp period exist
	mppPeriod, err := uc.MPPPeriodRepo.FindById(*req.MPPPeriodID)
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Update] error when find mpp period by id: %v", err)
		return nil, err
	}

	if mppPeriod == nil {
		uc.Log.Errorf("[MPRequestUseCase.Update] mpp period with id %s is not exist", req.MPPPeriodID.String())
		return nil, errors.New("mpp period is not exist")
	}

	// check portal data
	portalResponse, err := uc.MPRequestHelper.CheckPortalData(req)
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Update] error when check portal data: %v", err)
		return nil, err
	}

	mpRequestHeader, err := uc.MPRequestRepository.Update(dto.ConvertToEntity(req))
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Update] error when create mp request header: %v", err)
		return nil, err
	}

	// create request major
	for _, majorID := range req.MajorIDs {
		err := uc.RequestMajorRepository.DeleteByMPRequestHeaderID(mpRequestHeader.ID)
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.Update] error when delete request major by mp request header id: %v", err)
			return nil, err
		}
		reqMajor, err := uc.RequestMajorRepository.Create(&entity.RequestMajor{
			MPRequestHeaderID: mpRequestHeader.ID,
			MajorID:           majorID,
		})
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.Update] error when create request major: %v", err)
			return nil, err
		}

		mpRequestHeader.RequestMajors = append(mpRequestHeader.RequestMajors, *reqMajor)
	}

	mpRequestHeader.OrganizationName = portalResponse.OrganizationName
	mpRequestHeader.OrganizationLocationName = portalResponse.OrganizationLocationName
	mpRequestHeader.ForOrganizationName = portalResponse.ForOrganizationName
	mpRequestHeader.ForOrganizationLocation = portalResponse.ForOrganizationLocationName
	mpRequestHeader.ForOrganizationStructure = portalResponse.ForOrganizationStructureName
	mpRequestHeader.JobName = portalResponse.JobName
	mpRequestHeader.RequestorName = portalResponse.RequestorName
	mpRequestHeader.DepartmentHeadName = portalResponse.DepartmentHeadName
	mpRequestHeader.EmpOrganizationName = portalResponse.EmpOrganizationName
	mpRequestHeader.JobLevelName = portalResponse.JobLevelName
	mpRequestHeader.JobLevel = portalResponse.JobLevel
	mpRequestHeader.MPPPeriod = *mppPeriod

	return dto.ConvertToResponse(mpRequestHeader), nil
}

func (uc *MPRequestUseCase) FindAllPaginated(page int, pageSize int, search string, filter map[string]interface{}) (*response.MPRequestPaginatedResponse, error) {
	mpRequestHeaders, total, err := uc.MPRequestRepository.FindAllPaginated(page, pageSize, search, filter)
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when find all paginated mp request headers: %v", err)
		return nil, err
	}

	var mpRequestHeaderResponses []response.MPRequestHeaderResponse
	for _, mpRequestHeader := range mpRequestHeaders {
		// check portal data
		portalResponse, err := uc.MPRequestHelper.CheckPortalData(dto.ConvertEntityToRequest(&mpRequestHeader))
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when check portal data: %v", err)
			return nil, err
		}

		mpRequestHeader.OrganizationName = portalResponse.OrganizationName
		mpRequestHeader.OrganizationLocationName = portalResponse.OrganizationLocationName
		mpRequestHeader.ForOrganizationName = portalResponse.ForOrganizationName
		mpRequestHeader.ForOrganizationLocation = portalResponse.ForOrganizationLocationName
		mpRequestHeader.ForOrganizationStructure = portalResponse.ForOrganizationStructureName
		mpRequestHeader.JobName = portalResponse.JobName
		mpRequestHeader.RequestorName = portalResponse.RequestorName
		mpRequestHeader.DepartmentHeadName = portalResponse.DepartmentHeadName
		mpRequestHeader.VpGmDirectorName = portalResponse.VpGmDirectorName
		mpRequestHeader.CeoName = portalResponse.CeoName
		mpRequestHeader.HrdHoUnitName = portalResponse.HrdHoUnitName
		mpRequestHeader.EmpOrganizationName = portalResponse.EmpOrganizationName
		mpRequestHeader.JobLevelName = portalResponse.JobLevelName
		mpRequestHeader.JobLevel = portalResponse.JobLevel

		mpRequestHeaderResponses = append(mpRequestHeaderResponses, *dto.ConvertToResponse(&mpRequestHeader))
	}

	return &response.MPRequestPaginatedResponse{
		MPRequestHeader: mpRequestHeaderResponses,
		Total:           total,
	}, nil
}

func (uc *MPRequestUseCase) UpdateStatusHeader(req *request.UpdateMPRequestHeaderRequest) error {
	// check if mp request header is exist
	mpRequestHeader, err := uc.MPRequestRepository.FindById(uuid.MustParse(req.ID))
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.UpdateStatusHeader] error when find mp request header by id: %v", err)
		return err
	}

	if mpRequestHeader == nil {
		uc.Log.Errorf("[MPRequestUseCase.UpdateStatusHeader] mp request header with id %s is not exist", req.ID)
		return errors.New("mp request header is not exist")
	}

	// check if approver ID is exist
	approverExist, err := uc.EmpMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
		ID: req.ApproverID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.UpdateStatusHeader] error when send find employee by id message: %v", err)
		return err
	}
	if approverExist == nil {
		uc.Log.Errorf("[MPRequestUseCase.UpdateStatusHeader] approver with id %s is not exist", req.ApproverID.String())
		return errors.New("approver is not exist")
	}

	var approvalHistory *entity.MPRequestApprovalHistory

	if req.Status != entity.MPRequestStatusDraft && req.Status != entity.MPRequestStatusSubmitted && req.Status != entity.MPRequestStatusCompleted {
		approvalHistory = &entity.MPRequestApprovalHistory{
			MPRequestHeaderID: mpRequestHeader.ID,
			ApproverID:        req.ApproverID,
			ApproverName:      approverExist.Name,
			Level:             string(req.Level),
			Notes:             req.Notes,
			Status: func() entity.MPRequestApprovalHistoryStatus {
				if req.Status == entity.MPRequestStatusRejected {
					return entity.MPRequestApprovalHistoryStatusRejected
				} else if req.Status == entity.MPRequestStatusApproved {
					return entity.MPRequestApprovalHistoryStatusApproved
				} else if req.Status == entity.MPRequestStatusNeedApproval {
					return entity.MPRequestApprovalHistoryStatusNeedApproval
				}
				return entity.MPRequestApprovalHistoryStatusRejected
			}(),
		}
	}

	err = uc.MPRequestRepository.UpdateStatusHeader(uuid.MustParse(req.ID), string(req.Status), req.ApproverID.String(), approvalHistory)
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.UpdateStatusHeader] error when update mp request header: %v", err)
		return err
	}

	if req.Attachments != nil {
		for _, attachment := range req.Attachments {
			_, err := uc.MPRequestRepository.StoreAttachmentToApprovalHistory(approvalHistory, entity.ManpowerAttachment{
				FileName:  attachment.FileName,
				FileType:  attachment.FileType,
				FilePath:  attachment.FilePath,
				OwnerType: "mp_request_approval_histories",
				OwnerID:   approvalHistory.ID,
			})

			if err != nil {
				uc.Log.Errorf("[MPRequestUseCase.UpdateStatusHeader] error when store attachment to approval history: %v", err)
				return err
			}
		}
	}

	uc.Log.Infof("[MPRequestUseCase.UpdateStatusHeader] mp request header with id %s has been updated", string(req.ID))
	return nil
}

func MPRequestUseCaseFactory(viper *viper.Viper, log *logrus.Logger) IMPRequestUseCase {
	mpRequestRepository := repository.MPRequestRepositoryFactory(log)
	requestMajorRepository := repository.RequestMajorRepositoryFactory(log)
	organizationMessage := messaging.OrganizationMessageFactory(log)
	jobPlafonMessage := messaging.JobPlafonMessageFactory(log)
	userMessage := messaging.UserMessageFactory(log)
	mppPeriodRepo := repository.MPPPeriodRepositoryFactory(log)
	em := messaging.EmployeeMessageFactory(log)
	mprHelper := helper.MPRequestHelperFactory(log)
	return NewMPRequestUseCase(
		viper,
		log,
		mpRequestRepository,
		requestMajorRepository,
		organizationMessage,
		jobPlafonMessage,
		userMessage,
		mppPeriodRepo,
		em,
		mprHelper,
	)
}
