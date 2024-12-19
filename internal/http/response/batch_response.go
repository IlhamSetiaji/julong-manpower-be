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

type RealDocumentBatchResponse struct {
	Overall             DocumentBatchResponse         `json:"overall"`
	OrganizationOverall []OrganizationOverallResponse `json:"organization_overall"`
}

type OrganizationOverallResponse struct {
	Overall         DocumentBatchResponse   `json:"overall"`
	LocationOverall []DocumentBatchResponse `json:"location_overall"`
}

type DocumentBatchResponse struct {
	OperatingUnit string             `json:"operating_unit"`
	BudgetYear    string             `json:"budget_year"`
	Grade         GradeBatchResponse `json:"grade"`
}

type GradeBatchResponse struct {
	Executive    []DocumentCalculationBatchResponse `json:"executive"`
	NonExecutive []DocumentCalculationBatchResponse `json:"non_executive"`
	Total        []DocumentCalculationBatchResponse `json:"total"`
}

type DocumentCalculationBatchResponse struct {
	JobLevelName string `json:"job_level_name"`
	Existing     int    `json:"existing"`
	Promote      int    `json:"promote"`
	Recruit      int    `json:"recruit"`
	Total        int    `json:"total"`
	IsTotal      bool   `json:"is_total"`
}

type CompletedBatchResponse struct {
	ID             uuid.UUID        `json:"id"`
	DocumentNumber string           `json:"document_number"`
	DocumentDate   time.Time        `json:"document_date"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	MPPPeriod      MPPeriodResponse `json:"mpp_period"`
}
