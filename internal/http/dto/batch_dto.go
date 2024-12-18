package dto

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IBatchDTO interface {
	ConvertBatchHeaderEntityToResponse(batch *entity.BatchHeader) *response.BatchResponse
	ConvertToDocumentBatchResponse(batch *entity.BatchHeader, operatingUnit string) *response.DocumentBatchResponse
	ConvertToDocumentCalculationBatchResponse(mpPlanningLine entity.MPPlanningLine, isTotal bool) *response.DocumentCalculationBatchResponse
	ConvertDocumentCalculationBatchResponses(mpPlanningLines []entity.MPPlanningLine) []response.DocumentCalculationBatchResponse
	// ConvertRealDocumentBatchResponse(batch *entity.BatchHeader) *response.RealDocumentBatchResponse
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

// func (d *BatchDTO) ConvertRealDocumentBatchResponse(batch *entity.BatchHeader) *response.RealDocumentBatchResponse {
// 	return &response.RealDocumentBatchResponse{
// 		Overall: *d.ConvertToDocumentBatchResponse(batch, "Julong"),
// 		OrganizationOverall: func() []response.OrganizationOverallResponse {
// 			var organizationOverall []response.OrganizationOverallResponse
// 			for _, bl := range batch.BatchLines {

// 			return organizationOverall
// 		}(),
// 	}
// }

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

func (d *BatchDTO) ConvertToDocumentBatchResponse(batch *entity.BatchHeader, operatingUnit string) *response.DocumentBatchResponse {
	return &response.DocumentBatchResponse{
		OperatingUnit: operatingUnit,
		BudgetYear:    batch.DocumentDate.Format("2006"),
		Grade: response.GradeBatchResponse{
			Executive: func() []response.DocumentCalculationBatchResponse {
				var executive []response.DocumentCalculationBatchResponse
				for _, bl := range batch.BatchLines {
					for _, mpl := range bl.MPPlanningHeader.MPPlanningLines {
						if mpl.JobLevel > 3 {
							executive = append(executive, *d.ConvertToDocumentCalculationBatchResponse(mpl, false))
						}
					}
				}
				return executive
			}(),
			NonExecutive: func() []response.DocumentCalculationBatchResponse {
				var nonExecutive []response.DocumentCalculationBatchResponse
				for _, bl := range batch.BatchLines {
					for _, mpl := range bl.MPPlanningHeader.MPPlanningLines {
						if mpl.JobLevel <= 3 {
							nonExecutive = append(nonExecutive, *d.ConvertToDocumentCalculationBatchResponse(mpl, false))
						}
					}
				}
				return nonExecutive
			}(),
			Total: func() []response.DocumentCalculationBatchResponse {
				var total []response.DocumentCalculationBatchResponse
				var totalExisting, totalPromote, totalRecruit, totalOverall int
				for _, bl := range batch.BatchLines {
					for _, mpl := range bl.MPPlanningHeader.MPPlanningLines {
						totalExisting += mpl.Existing
						totalPromote += mpl.Promotion
						totalRecruit += mpl.RecruitPH + mpl.RecruitMT
						totalOverall += mpl.Existing + mpl.Promotion + mpl.RecruitPH + mpl.RecruitMT
					}
					total = append(total, response.DocumentCalculationBatchResponse{
						JobLevelName: "Total",
						Existing:     totalExisting,
						Promote:      totalPromote,
						Recruit:      totalRecruit,
						Total:        totalOverall,
						IsTotal:      true,
					})
				}
				return total
			}(),
		},
	}
}

func (d *BatchDTO) ConvertToDocumentCalculationBatchResponse(mpPlanningLine entity.MPPlanningLine, isTotal bool) *response.DocumentCalculationBatchResponse {
	return &response.DocumentCalculationBatchResponse{
		JobLevelName: mpPlanningLine.JobLevelName,
		Existing:     mpPlanningLine.Existing,
		Promote:      mpPlanningLine.Promotion,
		Recruit:      mpPlanningLine.RecruitPH + mpPlanningLine.RecruitMT,
		Total:        mpPlanningLine.Existing + mpPlanningLine.Promotion + mpPlanningLine.RecruitPH + mpPlanningLine.RecruitMT,
		IsTotal:      isTotal,
	}
}

func (d *BatchDTO) ConvertDocumentCalculationBatchResponses(mpPlanningLines []entity.MPPlanningLine) []response.DocumentCalculationBatchResponse {
	var responses []response.DocumentCalculationBatchResponse
	for _, mpPlanningLine := range mpPlanningLines {
		responses = append(responses, *d.ConvertToDocumentCalculationBatchResponse(mpPlanningLine, false))
	}
	return responses
}

func BatchDTOFactory(log *logrus.Logger) IBatchDTO {
	batchLineDTO := BatchLineDTOFactory(log)
	return NewBatchDTO(log, batchLineDTO)
}
