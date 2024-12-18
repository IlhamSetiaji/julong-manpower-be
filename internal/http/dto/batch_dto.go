package dto

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IBatchDTO interface {
	ConvertBatchHeaderEntityToResponse(batch *entity.BatchHeader) *response.BatchResponse
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
		Status:         string(batch.Status),
		CreatedAt:      batch.CreatedAt,
		UpdatedAt:      batch.UpdatedAt,
		BatchLines:     d.batchLineDTO.ConvertBatchLineEntitiesToResponse(&batch.BatchLines),
	}
}

func BatchDTOFactory(log *logrus.Logger) IBatchDTO {
	batchLineDTO := BatchLineDTOFactory(log)
	return NewBatchDTO(log, batchLineDTO)
}
