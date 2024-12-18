package response

import (
	"time"

	"github.com/google/uuid"
)

type BatchResponse struct {
	ID             uuid.UUID            `json:"id"`
	DocumentNumber string               `json:"document_number"`
	DocumentDate   time.Time            `json:"document_date"`
	ApproverID     *uuid.UUID           `json:"approver_id"`
	ApproverName   string               `json:"approver_name"`
	Status         string               `json:"status"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	BatchLines     []*BatchLineResponse `json:"batch_lines"`
}

type BatchLineResponse struct {
	ID                       uuid.UUID                 `json:"id"`
	BatchHeaderID            uuid.UUID                 `json:"batch_header_id"`
	MPPlanningHeaderID       uuid.UUID                 `json:"mp_planning_header_id"`
	OrganizationID           uuid.UUID                 `json:"organization_id"`
	OrganizationName         string                    `json:"organization_name"`
	OrganizationLocationID   uuid.UUID                 `json:"organization_location_id"`
	OrganizationLocationName string                    `json:"organization_location_name"`
	CreatedAt                time.Time                 `json:"created_at"`
	UpdatedAt                time.Time                 `json:"updated_at"`
	MPPlanningHeader         *MPPlanningHeaderResponse `json:"mp_planning_header"`
}
