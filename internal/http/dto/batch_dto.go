package dto

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
)

func ConvertBatchHeaderEntityToResponse(batch *entity.BatchHeader) *response.BatchResponse {
	return &response.BatchResponse{
		ID:             batch.ID,
		DocumentNumber: batch.DocumentNumber,
		Status:         string(batch.Status),
		CreatedAt:      batch.CreatedAt,
		UpdatedAt:      batch.UpdatedAt,
	}
}
