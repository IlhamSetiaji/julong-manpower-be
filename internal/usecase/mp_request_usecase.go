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

type IMPRequestUseCase interface {
	Create(req *request.CreateMPRequestHeaderRequest) (*response.MPRequestHeaderResponse, error)
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
	requestorExist, err := uc.EmpMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
		ID: req.RequestorID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] error when send find employee by id message: %v", err)
		return nil, err
	}

	if requestorExist == nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] requestor with id %s is not exist", req.RequestorID.String())
		return nil, errors.New("requestor is not exist")
	}

	// check if department head is exist
	var deptHeadExist *response.SendFindUserByIDResponse
	if req.DepartmentHead != nil {
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
	} else {
		deptHeadExist = &response.SendFindUserByIDResponse{}
	}

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

	// check if emp organization is exist
	empOrgExist, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: req.EmpOrganizationID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] error when send find emp organization by id message: %v", err)
		return nil, err
	}

	// check if job level is exist
	jobLevelExist, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
		ID: req.JobLevelID.String(),
	})
	if err != nil {
		uc.Log.Errorf("[MPRequestUseCase.Create] error when send find job level by id message: %v", err)
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

	mpRequestHeader.OrganizationName = orgExist.Name
	mpRequestHeader.OrganizationLocationName = orgLocExist.Name
	mpRequestHeader.ForOrganizationName = forOrgExist.Name
	mpRequestHeader.ForOrganizationLocation = forOrgLocExist.Name
	mpRequestHeader.ForOrganizationStructure = forOrgStructExist.Name
	mpRequestHeader.JobName = jobExist.Name
	mpRequestHeader.RequestorName = requestorExist.Name
	mpRequestHeader.DepartmentHeadName = deptHeadExist.Name
	mpRequestHeader.EmpOrganizationName = empOrgExist.Name
	mpRequestHeader.JobLevelName = jobLevelExist.Name
	mpRequestHeader.JobLevel = int(jobLevelExist.Level)
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
		// check if organization is exist
		orgExist, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
			ID: mpRequestHeader.OrganizationID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] error when send find organization by id message: %v", err)
			return nil, err
		}

		if orgExist == nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] organization with id %s is not exist", mpRequestHeader.OrganizationID.String())
			return nil, errors.New("organization is not exist")
		}

		// check if organization location is exist
		orgLocExist, err := uc.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
			ID: mpRequestHeader.OrganizationLocationID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] error when send find organization location by id message: %v", err)
			return nil, err
		}

		if orgLocExist == nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] organization location with id %s is not exist", mpRequestHeader.OrganizationLocationID.String())
			return nil, errors.New("organization location is not exist")
		}

		// check if for organization is exist
		forOrgExist, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
			ID: mpRequestHeader.ForOrganizationID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] error when send find for organization by id message: %v", err)
			return nil, err
		}

		if forOrgExist == nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] for organization with id %s is not exist", mpRequestHeader.ForOrganizationID.String())
			return nil, errors.New("for organization is not exist")
		}

		// check if for organization location is exist
		forOrgLocExist, err := uc.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
			ID: mpRequestHeader.ForOrganizationLocationID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] error when send find for organization location by id message: %v", err)
			return nil, err
		}

		if forOrgLocExist == nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] for organization location with id %s is not exist", mpRequestHeader.ForOrganizationLocationID.String())
			return nil, errors.New("for organization location is not exist")
		}

		// check if for organization structure is exist
		forOrgStructExist, err := uc.OrganizationMessage.SendFindOrganizationStructureByIDMessage(request.SendFindOrganizationStructureByIDMessageRequest{
			ID: mpRequestHeader.ForOrganizationStructureID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] error when send find for organization structure by id message: %v", err)
			return nil, err
		}

		if forOrgStructExist == nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] for organization structure with id %s is not exist", mpRequestHeader.ForOrganizationStructureID.String())
			return nil, errors.New("for organization structure is not exist")
		}

		// check if job ID is exist
		jobExist, err := uc.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
			ID: mpRequestHeader.JobID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] error when send find job by id message: %v", err)
			return nil, err
		}

		if jobExist == nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] job with id %s is not exist", mpRequestHeader.JobID.String())
			return nil, errors.New("job is not exist")
		}

		// check if requestor ID is exist
		requestorExist, err := uc.EmpMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: mpRequestHeader.RequestorID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] error when send find employee by id message: %v", err)
			return nil, err
		}

		if requestorExist == nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] requestor with id %s is not exist", mpRequestHeader.RequestorID.String())
			return nil, errors.New("requestor is not exist")
		}

		// check if department head is exist
		var deptHeadExist *response.EmployeeResponse
		if mpRequestHeader.DepartmentHead != nil {
			deptHeadExist, err = uc.EmpMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
				ID: mpRequestHeader.DepartmentHead.String(),
			})
			if err != nil {
				uc.Log.Errorf("[MPRequestUseCase.Create] error when send find user by id message: %v", err)
				return nil, err
			}

			if deptHeadExist == nil {
				uc.Log.Errorf("[MPRequestUseCase.Create] department head with id %s is not exist", mpRequestHeader.DepartmentHead.String())
				return nil, errors.New("department head is not exist")
			}
		} else {
			deptHeadExist = &response.EmployeeResponse{}
		}

		// check if vp gm director is exist
		var vpGmDirectorExist *response.EmployeeResponse
		if mpRequestHeader.VpGmDirector != nil {
			vpGmDirectorExist, err = uc.EmpMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
				ID: mpRequestHeader.VpGmDirector.String(),
			})
			if err != nil {
				uc.Log.Errorf("[MPRequestUseCase.Create] error when send find user by id message: %v", err)
				return nil, err
			}
		} else {
			vpGmDirectorExist = &response.EmployeeResponse{}
		}

		// check if ceo is exist
		var ceoExist *response.EmployeeResponse
		if mpRequestHeader.CEO != nil {
			ceoExist, err = uc.EmpMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
				ID: mpRequestHeader.CEO.String(),
			})
			if err != nil {
				uc.Log.Errorf("[MPRequestUseCase.Create] error when send find user by id message: %v", err)
				return nil, err
			}
		} else {
			ceoExist = &response.EmployeeResponse{}
		}

		// check if hrd ho unit is exist
		var hrdHoUnitExist *response.EmployeeResponse
		if mpRequestHeader.HrdHoUnit != nil {
			hrdHoUnitExist, err = uc.EmpMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
				ID: mpRequestHeader.HrdHoUnit.String(),
			})
			if err != nil {
				uc.Log.Errorf("[MPRequestUseCase.Create] error when send find user by id message: %v", err)
				return nil, err
			}
		} else {
			hrdHoUnitExist = &response.EmployeeResponse{}
		}

		// check if emp organization is exist
		empOrgExist, err := uc.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
			ID: mpRequestHeader.EmpOrganizationID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] error when send find emp organization by id message: %v", err)
			return nil, err
		}

		// check if job level is exist
		jobLevelExist, err := uc.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
			ID: mpRequestHeader.JobLevelID.String(),
		})
		if err != nil {
			uc.Log.Errorf("[MPRequestUseCase.Create] error when send find job level by id message: %v", err)
			return nil, err
		}

		mpRequestHeader.OrganizationName = orgExist.Name
		mpRequestHeader.OrganizationLocationName = orgLocExist.Name
		mpRequestHeader.ForOrganizationName = forOrgExist.Name
		mpRequestHeader.ForOrganizationLocation = forOrgLocExist.Name
		mpRequestHeader.ForOrganizationStructure = forOrgStructExist.Name
		mpRequestHeader.JobName = jobExist.Name
		mpRequestHeader.RequestorName = requestorExist.Name
		mpRequestHeader.DepartmentHeadName = deptHeadExist.Name
		mpRequestHeader.VpGmDirectorName = vpGmDirectorExist.Name
		mpRequestHeader.CeoName = ceoExist.Name
		mpRequestHeader.HrdHoUnitName = hrdHoUnitExist.Name
		mpRequestHeader.EmpOrganizationName = empOrgExist.Name
		mpRequestHeader.JobLevelName = jobLevelExist.Name
		mpRequestHeader.JobLevel = int(jobLevelExist.Level)

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
	return NewMPRequestUseCase(viper, log, mpRequestRepository, requestMajorRepository, organizationMessage, jobPlafonMessage, userMessage, mppPeriodRepo, em)
}
