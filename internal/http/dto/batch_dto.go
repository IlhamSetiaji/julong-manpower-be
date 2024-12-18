package dto

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IBatchDTO interface {
	ConvertBatchHeaderEntityToResponse(batch *entity.BatchHeader) *response.BatchResponse
	ConvertToDocumentBatchResponse(batch *entity.BatchHeader, operatingUnit string) *response.DocumentBatchResponse
	ConvertToDocumentCalculationBatchResponse(mpPlanningLine entity.MPPlanningLine) *response.DocumentCalculationBatchResponse
}

type BatchDTO struct {
	log          *logrus.Logger
	batchLineDTO IBatchLineDTO
}

func NewBatchDTO(log *logrus.Logger, batchLineDTO IBatchLineDTO) IBatchDTO {
	return &BatchDTO{
		log:          log,
		batchLineDTO: batchLineDTO,
	}
}

func (d *BatchDTO) ConvertBatchHeaderEntityToResponse(batch *entity.BatchHeader) *response.BatchResponse {
	return &response.BatchResponse{
		ID:             batch.ID,
		DocumentNumber: batch.DocumentNumber,
		DocumentDate:   batch.DocumentDate,
		Status:         string(batch.Status),
		CreatedAt:      batch.CreatedAt,
		UpdatedAt:      batch.UpdatedAt,
		BatchLines:     d.batchLineDTO.ConvertBatchLineEntitiesToResponse(&batch.BatchLines),
	}
}

// func (d *BatchDTO) ConvertToDocumentBatchResponse(batch *entity.BatchHeader, operatingUnit string) *response.DocumentBatchResponse {
// 	return &response.DocumentBatchResponse{
// 		OperatingUnit: operatingUnit,
// 		BudgetYear:    batch.DocumentDate.Format("2006"),
// 		Grade:
// 	}
// }

// func (d *BatchDTO) ConvertToDocumentCalculationBatchResponse(mpPlanningLine (entity.MPPlanningLine)) *response.DocumentCalculationBatchResponse {
// 	return

func BatchDTOFactory(log *logrus.Logger) IBatchDTO {
	batchLineDTO := BatchLineDTOFactory(log)
	return NewBatchDTO(log, batchLineDTO)
}
