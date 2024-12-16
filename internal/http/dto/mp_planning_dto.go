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
