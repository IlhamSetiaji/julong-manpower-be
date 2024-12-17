package dto

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/spf13/viper"
)

func ConvertMPPlanningApprovalHistoryToResponse(approvalHistories *entity.MPPlanningApprovalHistory, viper *viper.Viper) *response.MPPlanningApprovalHistoryResponse {
	return &response.MPPlanningApprovalHistoryResponse{
		ID:                 approvalHistories.ID,
		MPPlanningHeaderID: approvalHistories.MPPlanningHeaderID,
		ApproverID:         approvalHistories.ApproverID,
		ApproverName:       approvalHistories.ApproverName,
		Notes:              approvalHistories.Notes,
		Level:              approvalHistories.Level,
		Status:             approvalHistories.Status,
		Attachments:        ConvertManpowerAttachmentsToResponse(&approvalHistories.ManpowerAttachments, viper),
	}
}

func ConvertMPPlanningApprovalHistoriesToResponse(approvalHistories *[]entity.MPPlanningApprovalHistory, viper *viper.Viper) []*response.MPPlanningApprovalHistoryResponse {
	var response []*response.MPPlanningApprovalHistoryResponse
	for _, approvalHistory := range *approvalHistories {
		response = append(response, ConvertMPPlanningApprovalHistoryToResponse(&approvalHistory, viper))
	}
	return response
}

func ConvertMPPlanningHeaderEntityToResponse(mpPlanningHeader *entity.MPPlanningHeader) *response.MPPlanningHeaderResponse {
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
