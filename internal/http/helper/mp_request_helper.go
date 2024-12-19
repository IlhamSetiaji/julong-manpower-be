package helper

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IMPRequestHelper interface {
	CheckPortalData(req *request.CreateMPRequestHeaderRequest) (*response.CheckPortalDataMPRequestResponse, error)
}

type MPRequestHelper struct {
	Log                 *logrus.Logger
	OrganizationMessage messaging.IOrganizationMessage
	JobPlafonMessage    messaging.IJobPlafonMessage
	UserMessage         messaging.IUserMessage
	EmpMessage          messaging.IEmployeeMessage
}

func NewMPRequestHelper(
	log *logrus.Logger,
	organizationMessage messaging.IOrganizationMessage,
	jobPlafonMessage messaging.IJobPlafonMessage,
	userMessage messaging.IUserMessage,
	em messaging.IEmployeeMessage,
) IMPRequestHelper {
	return &MPRequestHelper{
		Log:                 log,
		OrganizationMessage: organizationMessage,
		JobPlafonMessage:    jobPlafonMessage,
		UserMessage:         userMessage,
		EmpMessage:          em,
	}
}

func MPRequestHelperFactory(log *logrus.Logger) IMPRequestHelper {
	organizationMessage := messaging.OrganizationMessageFactory(log)
	jobPlafonMessage := messaging.JobPlafonMessageFactory(log)
	userMessage := messaging.UserMessageFactory(log)
	em := messaging.EmployeeMessageFactory(log)
	return NewMPRequestHelper(log, organizationMessage, jobPlafonMessage, userMessage, em)
}

func (h *MPRequestHelper) CheckPortalData(req *request.CreateMPRequestHeaderRequest) (*response.CheckPortalDataMPRequestResponse, error) {
	// check if organization is exist
	orgExist, err := h.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: req.OrganizationID.String(),
	})
	if err != nil {
		h.Log.Errorf("[MPRequestHelper] error when send find organization by id message: %v", err)
		return nil, err
	}

	if orgExist == nil {
		h.Log.Errorf("[MPRequestHelper] organization with id %s is not exist", req.OrganizationID.String())
		return nil, errors.New("organization is not exist")
	}

	// check if organization location is exist
	orgLocExist, err := h.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
		ID: req.OrganizationLocationID.String(),
	})
	if err != nil {
		h.Log.Errorf("[MPRequestHelper] error when send find organization location by id message: %v", err)
		return nil, err
	}

	if orgLocExist == nil {
		h.Log.Errorf("[MPRequestHelper] organization location with id %s is not exist", req.OrganizationLocationID.String())
		return nil, errors.New("organization location is not exist")
	}

	// check if for organization is exist
	forOrgExist, err := h.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: req.ForOrganizationID.String(),
	})
	if err != nil {
		h.Log.Errorf("[MPRequestHelper] error when send find for organization by id message: %v", err)
		return nil, err
	}

	if forOrgExist == nil {
		h.Log.Errorf("[MPRequestHelper] for organization with id %s is not exist", req.ForOrganizationID.String())
		return nil, errors.New("for organization is not exist")
	}

	// check if for organization location is exist
	forOrgLocExist, err := h.OrganizationMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
		ID: req.ForOrganizationLocationID.String(),
	})
	if err != nil {
		h.Log.Errorf("[MPRequestHelper] error when send find for organization location by id message: %v", err)
		return nil, err
	}

	if forOrgLocExist == nil {
		h.Log.Errorf("[MPRequestHelper] for organization location with id %s is not exist", req.ForOrganizationLocationID.String())
		return nil, errors.New("for organization location is not exist")
	}

	// check if for organization structure is exist
	forOrgStructExist, err := h.OrganizationMessage.SendFindOrganizationStructureByIDMessage(request.SendFindOrganizationStructureByIDMessageRequest{
		ID: req.ForOrganizationStructureID.String(),
	})
	if err != nil {
		h.Log.Errorf("[MPRequestHelper] error when send find for organization structure by id message: %v", err)
		return nil, err
	}

	if forOrgStructExist == nil {
		h.Log.Errorf("[MPRequestHelper] for organization structure with id %s is not exist", req.ForOrganizationStructureID.String())
		return nil, errors.New("for organization structure is not exist")
	}

	// check if job ID is exist
	jobExist, err := h.JobPlafonMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
		ID: req.JobID.String(),
	})
	if err != nil {
		h.Log.Errorf("[MPRequestHelper] error when send find job by id message: %v", err)
		return nil, err
	}

	if jobExist == nil {
		h.Log.Errorf("[MPRequestHelper] job with id %s is not exist", req.JobID.String())
		return nil, errors.New("job is not exist")
	}

	// check if requestor ID is exist
	requestorExist, err := h.EmpMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
		ID: req.RequestorID.String(),
	})
	if err != nil {
		h.Log.Errorf("[MPRequestHelper] error when send find employee by id message: %v", err)
		return nil, err
	}

	if requestorExist == nil {
		h.Log.Errorf("[MPRequestHelper] requestor with id %s is not exist", req.RequestorID.String())
		return nil, errors.New("requestor is not exist")
	}

	// check if department head is exist
	var deptHeadExist *response.EmployeeResponse
	if req.DepartmentHead != nil {
		deptHeadExist, err = h.EmpMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: req.DepartmentHead.String(),
		})
		if err != nil {
			h.Log.Errorf("[MPRequestHelper] error when send find user by id message: %v", err)
			return nil, err
		}

		if deptHeadExist == nil {
			h.Log.Errorf("[MPRequestHelper] department head with id %s is not exist", req.DepartmentHead.String())
			return nil, errors.New("department head is not exist")
		}
	} else {
		deptHeadExist = &response.EmployeeResponse{}
	}

	// check if vp gm director is exist
	var vpGmDirectorExist *response.EmployeeResponse
	if req.VpGmDirector != nil {
		vpGmDirectorExist, err = h.EmpMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: req.VpGmDirector.String(),
		})
		if err != nil {
			h.Log.Errorf("[MPRequestHelper] error when send find user by id message: %v", err)
			return nil, err
		}
	} else {
		vpGmDirectorExist = &response.EmployeeResponse{}
	}

	// check if ceo is exist
	var ceoExist *response.EmployeeResponse
	if req.CEO != nil {
		ceoExist, err = h.EmpMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: req.CEO.String(),
		})
		if err != nil {
			h.Log.Errorf("[MPRequestHelper] error when send find user by id message: %v", err)
			return nil, err
		}

		if ceoExist == nil {
			h.Log.Errorf("[MPRequestHelper] ceo with id %s is not exist", req.CEO.String())
			return nil, errors.New("ceo is not exist")
		}
	} else {
		ceoExist = &response.EmployeeResponse{}
	}

	// check if hrd ho unit is exist
	var hrdHoUnitExist *response.EmployeeResponse
	if req.HrdHoUnit != nil {
		hrdHoUnitExist, err = h.EmpMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: req.HrdHoUnit.String(),
		})
		if err != nil {
			h.Log.Errorf("[MPRequestHelper] error when send find user by id message: %v", err)
			return nil, err
		}
	} else {
		hrdHoUnitExist = &response.EmployeeResponse{}
	}

	// check if emp organization is exist
	empOrgExist, err := h.OrganizationMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
		ID: req.EmpOrganizationID.String(),
	})
	if err != nil {
		h.Log.Errorf("[MPRequestHelper] error when send find emp organization by id message: %v", err)
		return nil, err
	}

	// check if job level is exist
	jobLevelExist, err := h.JobPlafonMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
		ID: req.JobLevelID.String(),
	})
	if err != nil {
		h.Log.Errorf("[MPRequestHelper] error when send find job level by id message: %v", err)
		return nil, err
	}

	if jobLevelExist == nil {
		h.Log.Errorf("[MPRequestHelper] job level with id %s is not exist", req.JobLevelID.String())
		return nil, errors.New("job level is not exist")
	}

	return &response.CheckPortalDataMPRequestResponse{
		OrganizationName:             orgExist.Name,
		OrganizationCategory:         orgExist.OrganizationCategory,
		OrganizationLocationName:     orgLocExist.Name,
		ForOrganizationName:          forOrgExist.Name,
		ForOrganizationLocationName:  forOrgLocExist.Name,
		ForOrganizationStructureName: forOrgStructExist.Name,
		JobName:                      jobExist.Name,
		RequestorName:                requestorExist.Name,
		DepartmentHeadName:           deptHeadExist.Name,
		VpGmDirectorName:             vpGmDirectorExist.Name,
		CeoName:                      ceoExist.Name,
		HrdHoUnitName:                hrdHoUnitExist.Name,
		EmpOrganizationName:          empOrgExist.Name,
		JobLevelName:                 jobLevelExist.Name,
		JobLevel:                     int(jobLevelExist.Level),
	}, nil
}
