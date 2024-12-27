package dto

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IMPPlanningDTO interface {
	ConvertMPPlanningApprovalHistoryToResponse(approvalHistories *entity.MPPlanningApprovalHistory, viper *viper.Viper) *response.MPPlanningApprovalHistoryResponse
	ConvertMPPlanningApprovalHistoriesToResponse(approvalHistories *[]entity.MPPlanningApprovalHistory, viper *viper.Viper) []*response.MPPlanningApprovalHistoryResponse
	ConvertMPPlanningHeaderEntityToResponse(mpPlanningHeader *entity.MPPlanningHeader) *response.MPPlanningHeaderResponse
	ConvertMPPlanningHeaderEntititesToResponse(mpPlanningHeaders *[]entity.MPPlanningHeader) []*response.MPPlanningHeaderResponse
}

type MPPlanningDTO struct {
	log        *logrus.Logger
	orgMessage messaging.IOrganizationMessage
	jobMessage messaging.IJobMessage
	jpMessage  messaging.IJobPlafonMessage
	empMessage messaging.IEmployeeMessage
}

func NewMPPlanningDTO(log *logrus.Logger, orgMessage messaging.IOrganizationMessage, jobMessage messaging.IJobMessage, jpMessage messaging.IJobPlafonMessage, empMessage messaging.IEmployeeMessage) IMPPlanningDTO {
	return &MPPlanningDTO{
		log:        log,
		orgMessage: orgMessage,
		jobMessage: jobMessage,
		jpMessage:  jpMessage,
		empMessage: empMessage,
	}
}

func (d *MPPlanningDTO) ConvertMPPlanningApprovalHistoryToResponse(approvalHistories *entity.MPPlanningApprovalHistory, viper *viper.Viper) *response.MPPlanningApprovalHistoryResponse {
	return &response.MPPlanningApprovalHistoryResponse{
		ID:                 approvalHistories.ID,
		MPPlanningHeaderID: approvalHistories.MPPlanningHeaderID,
		ApproverID:         approvalHistories.ApproverID,
		ApproverName:       approvalHistories.ApproverName,
		Notes:              approvalHistories.Notes,
		Level:              approvalHistories.Level,
		Status:             approvalHistories.Status,
		CreatedAt:          approvalHistories.CreatedAt,
		UpdatedAt:          approvalHistories.UpdatedAt,
		Attachments:        ConvertManpowerAttachmentsToResponse(&approvalHistories.ManpowerAttachments, viper),
	}
}

func (d *MPPlanningDTO) ConvertMPPlanningApprovalHistoriesToResponse(approvalHistories *[]entity.MPPlanningApprovalHistory, viper *viper.Viper) []*response.MPPlanningApprovalHistoryResponse {
	var response []*response.MPPlanningApprovalHistoryResponse
	for _, approvalHistory := range *approvalHistories {
		response = append(response, d.ConvertMPPlanningApprovalHistoryToResponse(&approvalHistory, viper))
	}
	return response
}

func (d *MPPlanningDTO) ConvertMPPlanningHeaderEntityToResponse(mpPlanningHeader *entity.MPPlanningHeader) *response.MPPlanningHeaderResponse {
	if mpPlanningHeader.OrganizationName == "" {
		d.log.Infof("Organization ID: %s", mpPlanningHeader.OrganizationID)
		organization, err := d.orgMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
			ID: mpPlanningHeader.OrganizationID.String(),
		})
		if err != nil {
			d.log.Errorf("[MPPlanningDTO.ConvertMPPlanningHeaderEntityToResponse] " + err.Error())
			mpPlanningHeader.OrganizationName = ""
		} else {
			mpPlanningHeader.OrganizationName = organization.Name
		}
	}

	if mpPlanningHeader.EmpOrganizationName == "" {
		organization, err := d.orgMessage.SendFindOrganizationByIDMessage(request.SendFindOrganizationByIDMessageRequest{
			ID: mpPlanningHeader.EmpOrganizationID.String(),
		})
		if err != nil {
			d.log.Errorf("[MPPlanningDTO.ConvertMPPlanningHeaderEntityToResponse] " + err.Error())
			mpPlanningHeader.EmpOrganizationName = ""
		} else {
			mpPlanningHeader.EmpOrganizationName = organization.Name
		}
	}

	if mpPlanningHeader.JobName == "" {
		job, err := d.jpMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
			ID: mpPlanningHeader.JobID.String(),
		})
		if err != nil {
			d.log.Errorf("[MPPlanningDTO.ConvertMPPlanningHeaderEntityToResponse] " + err.Error())
			mpPlanningHeader.JobName = ""
		} else {
			mpPlanningHeader.JobName = job.Name
		}
	}

	if mpPlanningHeader.RequestorName == "" || mpPlanningHeader.RequestorID == nil {
		employee, err := d.empMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: mpPlanningHeader.RequestorID.String(),
		})
		if err != nil {
			d.log.Errorf("[MPPlanningDTO.ConvertMPPlanningHeaderEntityToResponse] " + err.Error())
			mpPlanningHeader.RequestorName = ""
		} else {
			mpPlanningHeader.RequestorName = employee.Name
		}
	}

	if mpPlanningHeader.OrganizationLocationName == "" {
		organizationLocation, err := d.orgMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
			ID: mpPlanningHeader.OrganizationLocationID.String(),
		})
		if err != nil {
			d.log.Errorf("[MPPlanningDTO.ConvertMPPlanningHeaderEntityToResponse] " + err.Error())
			mpPlanningHeader.OrganizationLocationName = ""
		} else {
			mpPlanningHeader.OrganizationLocationName = organizationLocation.Name
		}
	}

	return &response.MPPlanningHeaderResponse{
		ID:                       mpPlanningHeader.ID,
		MPPPeriodID:              mpPlanningHeader.MPPPeriodID,
		OrganizationID:           mpPlanningHeader.OrganizationID,
		EmpOrganizationID:        mpPlanningHeader.EmpOrganizationID,
		JobID:                    mpPlanningHeader.JobID,
		DocumentNumber:           mpPlanningHeader.DocumentNumber,
		DocumentDate:             mpPlanningHeader.DocumentDate,
		Notes:                    mpPlanningHeader.Notes,
		TotalRecruit:             mpPlanningHeader.TotalRecruit,
		TotalPromote:             mpPlanningHeader.TotalPromote,
		Status:                   mpPlanningHeader.Status,
		RecommendedBy:            mpPlanningHeader.RecommendedBy,
		ApprovedBy:               mpPlanningHeader.ApprovedBy,
		ApproverManagerID:        mpPlanningHeader.ApproverManagerID,
		ApproverRecruitmentID:    mpPlanningHeader.ApproverRecruitmentID,
		RequestorID:              mpPlanningHeader.RequestorID,
		NotesAttach:              mpPlanningHeader.NotesAttach,
		OrganizationName:         mpPlanningHeader.OrganizationName,
		EmpOrganizationName:      mpPlanningHeader.EmpOrganizationName,
		JobName:                  mpPlanningHeader.JobName,
		RequestorName:            mpPlanningHeader.RequestorName,
		OrganizationLocationID:   mpPlanningHeader.OrganizationLocationID,
		OrganizationLocationName: mpPlanningHeader.OrganizationLocationName,
		CreatedAt:                mpPlanningHeader.CreatedAt,
		UpdatedAt:                mpPlanningHeader.UpdatedAt,
		MPPPeriod: &response.MPPeriodResponse{
			ID:              mpPlanningHeader.MPPPeriod.ID,
			Title:           mpPlanningHeader.MPPPeriod.Title,
			StartDate:       mpPlanningHeader.MPPPeriod.StartDate.Format("2006-01-02"),
			EndDate:         mpPlanningHeader.MPPPeriod.EndDate.Format("2006-01-02"),
			BudgetStartDate: mpPlanningHeader.MPPPeriod.BudgetStartDate.Format("2006-01-02"),
			BudgetEndDate:   mpPlanningHeader.MPPPeriod.BudgetEndDate.Format("2006-01-02"),
			Status:          mpPlanningHeader.MPPPeriod.Status,
			CreatedAt:       mpPlanningHeader.MPPPeriod.CreatedAt,
			UpdatedAt:       mpPlanningHeader.MPPPeriod.UpdatedAt,
		},
		MPPlanningLines: func() []*response.MPPlanningLineResponse {
			var lines []*response.MPPlanningLineResponse
			for _, line := range mpPlanningHeader.MPPlanningLines {
				if line.JobName == "" {
					job, err := d.jpMessage.SendFindJobByIDMessage(request.SendFindJobByIDMessageRequest{
						ID: line.JobID.String(),
					})
					if err != nil {
						d.log.Errorf("[MPPlanningDTO.ConvertMPPlanningHeaderEntityToResponse] " + err.Error())
					}
					line.JobName = job.Name
				}

				if line.JobLevelName == "" {
					jobLevel, err := d.jpMessage.SendFindJobLevelByIDMessage(request.SendFindJobLevelByIDMessageRequest{
						ID: line.JobLevelID.String(),
					})
					if err != nil {
						d.log.Errorf("[MPPlanningDTO.ConvertMPPlanningHeaderEntityToResponse] " + err.Error())
					}
					line.JobLevelName = jobLevel.Name
				}

				if line.OrganizationLocationName == "" {
					organizationLocation, err := d.orgMessage.SendFindOrganizationLocationByIDMessage(request.SendFindOrganizationLocationByIDMessageRequest{
						ID: line.OrganizationLocationID.String(),
					})
					if err != nil {
						d.log.Errorf("[MPPlanningDTO.ConvertMPPlanningHeaderEntityToResponse] " + err.Error())
					}
					line.OrganizationLocationName = organizationLocation.Name
				}

				lines = append(lines, &response.MPPlanningLineResponse{
					ID:                       line.ID,
					MPPlanningHeaderID:       line.MPPlanningHeaderID,
					OrganizationLocationID:   *line.OrganizationLocationID,
					JobLevelID:               *line.JobLevelID,
					JobID:                    *line.JobID,
					Existing:                 line.Existing,
					Recruit:                  line.Recruit,
					SuggestedRecruit:         line.SuggestedRecruit,
					Promotion:                line.Promotion,
					Total:                    line.Total,
					RemainingBalancePH:       line.RemainingBalancePH,
					RemainingBalanceMT:       line.RemainingBalanceMT,
					RecruitPH:                line.RecruitPH,
					RecruitMT:                line.RecruitMT,
					OrganizationLocationName: line.OrganizationLocationName,
					JobLevelName:             line.JobLevelName,
					JobName:                  line.JobName,
				})
			}
			return lines
		}(),
	}
}

func (d *MPPlanningDTO) ConvertMPPlanningHeaderEntititesToResponse(mpPlanningHeaders *[]entity.MPPlanningHeader) []*response.MPPlanningHeaderResponse {
	var response []*response.MPPlanningHeaderResponse
	for _, mpPlanningHeader := range *mpPlanningHeaders {
		response = append(response, d.ConvertMPPlanningHeaderEntityToResponse(&mpPlanningHeader))
	}
	return response
}

func MPPlanningDTOFactory(log *logrus.Logger) IMPPlanningDTO {
	orgMessage := messaging.OrganizationMessageFactory(log)
	jobMessage := messaging.JobMessageFactory(log)
	jpMessage := messaging.JobPlafonMessageFactory(log)
	empMessage := messaging.EmployeeMessageFactory(log)
	return NewMPPlanningDTO(log, orgMessage, jobMessage, jpMessage, empMessage)
}
