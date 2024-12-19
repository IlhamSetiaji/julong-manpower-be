package usecase

import (
	"errors"
	"fmt"
	"time"

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
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*response.MPRequestHeaderResponse, error)
	FindAllPaginated(page int, pageSize int, search string, filter map[string]interface{}) (*response.MPRequestPaginatedResponse, error)
	UpdateStatusHeader(req *request.UpdateMPRequestHeaderRequest) error
	GenerateDocumentNumber(dateNow time.Time) (string, error)
	CountTotalApprovalHistoryByStatus(headerID uuid.UUID, status entity.MPRequestApprovalHistoryStatus) (int64, error)
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
	MPRequestDTO           dto.IMPRequestDTO
	MPPlanningRepository   repository.IMPPlanningRepository
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
	mprDTO dto.IMPRequestDTO,
	mpPlanningRepository repository.IMPPlanningRepository,
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
		MPRequestDTO:           mprDTO,
		MPPlanningRepository:   mpPlanningRepository,
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

	mpRequestHeader, err := uc.MPRequestRepository.Create(uc.MPRequestDTO.ConvertToEntity(req))
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
	mpRequestHeader.OrganizationCategory = portalResponse.OrganizationCategory
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

	return uc.MPRequestDTO.ConvertToResponse(mpRequestHeader), nil
}

func (uc *MPRequestUseCase) CountTotalApprovalHistoryByStatus(headerID uuid.UUID, status entity.MPRequestApprovalHistoryStatus) (int64, error) {
	exist, err := uc.MPRequestRepository.FindById(headerID)
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.CountTotalApprovalHistoryByStatus] error when find mp request header by id: %v", err)
		return 0, err
	}

	if exist == nil {
		uc.Log.Errorf("[MPRequestUseCase.CountTotalApprovalHistoryByStatus] mp request header with id %s is not exist", headerID.String())
		return 0, errors.New("mp request header is not exist")
	}

	total, err := uc.MPRequestRepository.CountTotalApprovalHistoryByStatus(headerID, status)
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.CountTotalApprovalHistoryByStatus] error when count total approval history by status: %v", err)
		return 0, err
	}

	return total, nil
}

func (uc *MPRequestUseCase) Delete(id uuid.UUID) error {
	mpRequestHeader, err := uc.MPRequestRepository.FindById(id)
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Delete] error when find mp request header by id: %v", err)
		return err
	}

	if mpRequestHeader == nil {
		uc.Log.Errorf("[MPRequestUseCase.Delete] mp request header with id %s is not exist", id.String())
		return errors.New("mp request header is not exist")
	}

	err = uc.MPRequestRepository.DeleteHeader(id)
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Delete] error when delete mp request header: %v", err)
		return err
	}

	uc.Log.Infof("[MPRequestUseCase.Delete] mp request header with id %s has been deleted", id.String())
	return nil
}

func (uc *MPRequestUseCase) GenerateDocumentNumber(dateNow time.Time) (string, error) {
	foundMpRequestHeader, err := uc.MPRequestRepository.GetHeadersByDocumentDate(dateNow.Format("2006-01-02"))
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.GenerateDocumentNumber] error when get headers by document date: %v", err)
		return "", err
	}

	if len(foundMpRequestHeader) == 0 || foundMpRequestHeader == nil {
		return "MPR/" + dateNow.Format("20060102") + "/001", nil
	}

	return "MPR/" + dateNow.Format("20060102") + "/" + fmt.Sprintf("%03d", len(*&foundMpRequestHeader)+1), nil
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

	mpRequestHeader, err := uc.MPRequestRepository.Update(uc.MPRequestDTO.ConvertToEntity(req))
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
	mpRequestHeader.OrganizationCategory = portalResponse.OrganizationCategory
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

	return uc.MPRequestDTO.ConvertToResponse(mpRequestHeader), nil
}

func (uc *MPRequestUseCase) FindByID(id uuid.UUID) (*response.MPRequestHeaderResponse, error) {
	mpRequestHeader, err := uc.MPRequestRepository.FindById(id)
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.FindByID] error when find mp request header by id: %v", err)
		return nil, err
	}

	if mpRequestHeader == nil {
		uc.Log.Errorf("[MPRequestUseCase.FindByID] mp request header with id %s is not exist", id.String())
		return nil, errors.New("mp request header is not exist")
	}

	// check portal data
	portalResponse, err := uc.MPRequestHelper.CheckPortalData(uc.MPRequestDTO.ConvertEntityToRequest(mpRequestHeader))
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.FindByID] error when check portal data: %v", err)
		return nil, err
	}

	mpRequestHeader.OrganizationName = portalResponse.OrganizationName
	mpRequestHeader.OrganizationCategory = portalResponse.OrganizationCategory
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

	return uc.MPRequestDTO.ConvertToResponse(mpRequestHeader), nil
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
		portalResponse, err := uc.MPRequestHelper.CheckPortalData(uc.MPRequestDTO.ConvertEntityToRequest(&mpRequestHeader))
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.FindAllPaginated] error when check portal data: %v", err)
			return nil, err
		}

		mpRequestHeader.OrganizationName = portalResponse.OrganizationName
		mpRequestHeader.OrganizationCategory = portalResponse.OrganizationCategory
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

		mpRequestHeaderResponses = append(mpRequestHeaderResponses, *uc.MPRequestDTO.ConvertToResponse(&mpRequestHeader))
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

	if req.Level == entity.MPPRequestApprovalHistoryLevelHRDHO && req.Status == entity.MPRequestStatusCompleted {
		// find mp planning line by header id and job id
		// mpPlanningLine, err := uc.MPPlanningRepository.FindLineByHeaderIDAndJobID(*mpRequestHeader.MPPlanningHeaderID, *mpRequestHeader.JobID)
		// if err != nil {
		// 	uc.Log.Errorf("[MPRequestUseCase.UpdateStatusHeader] error when find mp planning line by header id and job id: %v", err)
		// 	return err
		// }

		// if mpPlanningLine == nil {
		// 	uc.Log.Errorf("[MPRequestUseCase.UpdateStatusHeader] mp planning line with header id %s and job id %s is not exist", mpRequestHeader.ID.String(), mpRequestHeader.JobID.String())
		// 	// return errors.New("mp planning line is not exist")
		// 	// return nil
		// }

		// if mpRequestHeader.RecruitmentType == entity.RecruitmentTypeEnumMT {
		// 	mpPlanningLine.RemainingBalanceMT = mpPlanningLine.RemainingBalanceMT - mpRequestHeader.TotalNeeds
		// } else if mpRequestHeader.RecruitmentType == entity.RecruitmentTypeEnumPH {
		// 	mpPlanningLine.RemainingBalancePH = mpPlanningLine.RemainingBalancePH - mpRequestHeader.TotalNeeds
		// }
		// _, err = uc.MPPlanningRepository.UpdateLineByHeaderIDAndJobID(mpRequestHeader.ID, *mpRequestHeader.JobID, mpPlanningLine)
		// if err != nil {
		// 	uc.Log.Errorf("[MPRequestUseCase.UpdateStatusHeader] error when update mp planning line by header id and job id: %v", err)
		// 	return err
		// }
		mpPlanningLines, err := uc.MPPlanningRepository.GetLinesByHeaderAndJobID(*mpRequestHeader.MPPlanningHeaderID, *mpRequestHeader.MPPlanningHeader.JobID)
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.UpdateStatusHeader] error when get mp planning lines by header and job id: %v", err)
			return err
		}

		if mpPlanningLines == nil {
			uc.Log.Errorf("[MPRequestUseCase.UpdateStatusHeader] mp planning lines with header id %s and job id %s is not exist", mpRequestHeader.MPPlanningHeaderID.String(), mpRequestHeader.JobID.String())
			return errors.New("mp planning lines is not exist")
		}

		for _, mpPlanningLine := range *mpPlanningLines {
			if mpRequestHeader.RecruitmentType == entity.RecruitmentTypeEnumMT {
				mpPlanningLine.RemainingBalanceMT = mpPlanningLine.RemainingBalanceMT - mpRequestHeader.TotalNeeds
			} else if mpRequestHeader.RecruitmentType == entity.RecruitmentTypeEnumPH {
				mpPlanningLine.RemainingBalancePH = mpPlanningLine.RemainingBalancePH - mpRequestHeader.TotalNeeds
			}
			_, err = uc.MPPlanningRepository.UpdateLineByHeaderIDAndJobID(*mpRequestHeader.MPPlanningHeaderID, *mpRequestHeader.JobID, &mpPlanningLine)
			if err != nil {
				uc.Log.Errorf("[MPRequestUseCase.UpdateStatusHeader] error when update mp planning line by header id and job id: %v", err)
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
	mprDTO := dto.MPRequestDTOFactory(log)
	mpPlanningRepo := repository.MPPlanningRepositoryFactory(log)
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
		mprDTO,
		mpPlanningRepo,
	)
}
