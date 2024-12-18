package request

import "github.com/IlhamSetiaji/julong-manpower-be/internal/entity"

type CreateBatchHeaderAndLinesRequest struct {
	DocumentNumber string                           `json:"document_number" validate:"omitempty,max=255"` // max length 255
	Status         entity.BatchHeaderApprovalStatus `json:"status" validate:"omitempty,BatchHeaderApprovalStatusValidation"`
	BatchLines     []struct {
		MPPlanningHeaderID     string `json:"mp_planning_header_id" validate:"required"`
		OrganizationID         string `json:"organization_id" validate:"required"`
		OrganizationLocationID string `json:"organization_location_id" validate:"required"`
	} `json:"batch_lines" validate:"required,dive"`
}
