package dto

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
)

func ConvertBatchLineEntityToResponse(batch *entity.BatchLine) *response.BatchLineResponse {
	return &response.BatchLineResponse{
		ID:                       batch.ID,
		MPPlanningHeaderID:       batch.MPPlanningHeaderID,
		OrganizationID:           *batch.OrganizationID,
		OrganizationLocationID:   *batch.OrganizationLocationID,
		OrganizationName:         batch.OrganizationName,
		OrganizationLocationName: batch.OrganizationLocationName,
		CreatedAt:                batch.CreatedAt,
		UpdatedAt:                batch.UpdatedAt,
		MPPlanningHeader:         ConvertMPPlanningHeaderEntityToResponse(&batch.MPPlanningHeader),
	}
}

func ConvertBatchLineEntitiesToResponse(batches *[]entity.BatchLine) []*response.BatchLineResponse {
	var response []*response.BatchLineResponse
	for _, batch := range *batches {
		response = append(response, ConvertBatchLineEntityToResponse(&batch))
	}
	return response
}
