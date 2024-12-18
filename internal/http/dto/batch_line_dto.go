package dto

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IBatchLineDTO interface {
	ConvertBatchLineEntityToResponse(batch *entity.BatchLine) *response.BatchLineResponse
	ConvertBatchLineEntitiesToResponse(batches *[]entity.BatchLine) []*response.BatchLineResponse
}

type BatchLineDTO struct {
	log         *logrus.Logger
	mpHeaderDTO IMPPlanningDTO
}

func NewBatchLineDTO(log *logrus.Logger, mpHeaderDTO IMPPlanningDTO) IBatchLineDTO {
	return &BatchLineDTO{
		log:         log,
		mpHeaderDTO: mpHeaderDTO,
	}
}

func (d *BatchLineDTO) ConvertBatchLineEntityToResponse(batch *entity.BatchLine) *response.BatchLineResponse {
	return &response.BatchLineResponse{
		ID:                       batch.ID,
		BatchHeaderID:            batch.BatchHeaderID,
		MPPlanningHeaderID:       batch.MPPlanningHeaderID,
		OrganizationID:           *batch.OrganizationID,
		OrganizationLocationID:   *batch.OrganizationLocationID,
		OrganizationName:         batch.OrganizationName,
		OrganizationLocationName: batch.OrganizationLocationName,
		CreatedAt:                batch.CreatedAt,
		UpdatedAt:                batch.UpdatedAt,
		MPPlanningHeader:         d.mpHeaderDTO.ConvertMPPlanningHeaderEntityToResponse(&batch.MPPlanningHeader),
	}
}

func (d *BatchLineDTO) ConvertBatchLineEntitiesToResponse(batches *[]entity.BatchLine) []*response.BatchLineResponse {
	var response []*response.BatchLineResponse
	for _, batch := range *batches {
		response = append(response, d.ConvertBatchLineEntityToResponse(&batch))
	}
	return response
}

func BatchLineDTOFactory(log *logrus.Logger) IBatchLineDTO {
	mpHeaderDTO := MPPlanningDTOFactory(log)
	return NewBatchLineDTO(log, mpHeaderDTO)
}
