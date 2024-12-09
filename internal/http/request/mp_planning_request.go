package request

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
)

type FindAllHeadersPaginatedMPPlanningRequest struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"pageSize" binding:"required"`
	Search   string `json:"search"`
}

type FindHeaderByIdMPPlanningRequest struct {
	ID string `json:"id" validate:"required"`
}

type ManpowerAttachmentRequest struct {
	FileName string `json:"file_name" validate:"required"`
	FilePath string `json:"file_path" validate:"required"`
	FileType string `json:"file_type" validate:"required"`
}

type CreateHeaderMPPlanningRequest struct {
	MPPPeriodID       uuid.UUID                   `json:"mpp_period_id" validate:"required"`
	OrganizationID    uuid.UUID                   `json:"organization_id" validate:"required"`
	EmpOrganizationID uuid.UUID                   `json:"emp_organization_id" validate:"required"`
	JobID             uuid.UUID                   `json:"job_id" validate:"required"` // job_id
	DocumentNumber    string                      `json:"document_number" validate:"required"`
	DocumentDate      string                      `json:"document_date" validate:"required,datetime=2006-01-02"`
	Notes             string                      `json:"notes" validate:"omitempty"`
	TotalRecruit      float64                     `json:"total_recruit" validate:"required"`
	TotalPromote      float64                     `json:"total_promote" validate:"required"`
	Status            entity.MPPlaningStatus      `json:"status" validate:"required,MPPlaningStatusValidation"`
	RecommendedBy     string                      `json:"recommended_by" validate:"required"`
	ApprovedBy        string                      `json:"approved_by" validate:"required"`
	RequestorID       uuid.UUID                   `json:"requestor_id" validate:"required"`
	NotesAttach       string                      `json:"notes_attach" validate:"omitempty"`
	Attachments       []ManpowerAttachmentRequest `json:"attachments" validate:"omitempty,dive"`
}

type UpdateHeaderMPPlanningRequest struct {
	ID                uuid.UUID                   `json:"id" validate:"required"`
	MPPPeriodID       uuid.UUID                   `json:"mpp_period_id" validate:"required"`
	OrganizationID    uuid.UUID                   `json:"organization_id" validate:"required"`
	EmpOrganizationID uuid.UUID                   `json:"emp_organization_id" validate:"required"`
	JobID             uuid.UUID                   `json:"job_id" validate:"required"` // job_id
	DocumentNumber    string                      `json:"document_number" validate:"required"`
	DocumentDate      string                      `json:"document_date" validate:"required,datetime=2006-01-02"`
	Notes             string                      `json:"notes" validate:"omitempty"`
	TotalRecruit      float64                     `json:"total_recruit" validate:"required"`
	TotalPromote      float64                     `json:"total_promote" validate:"required"`
	Status            entity.MPPlaningStatus      `json:"status" validate:"required,MPPlaningStatusValidation"`
	RecommendedBy     string                      `json:"recommended_by" validate:"required"`
	ApprovedBy        string                      `json:"approved_by" validate:"required"`
	RequestorID       uuid.UUID                   `json:"requestor_id" validate:"required"`
	NotesAttach       string                      `json:"notes_attach" validate:"omitempty"`
	Attachments       []ManpowerAttachmentRequest `json:"attachments" validate:"omitempty,dive"`
}

type DeleteHeaderMPPlanningRequest struct {
	ID string `json:"id" validate:"required"`
}

type FindAllLinesByHeaderIdPaginatedMPPlanningLineRequest struct {
	HeaderID string `json:"header_id" validate:"required"`
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"pageSize" binding:"required"`
	Search   string `json:"search"`
}

type FindLineByIdMPPlanningLineRequest struct {
	ID string `json:"id" validate:"required"`
}

type CreateLineMPPlanningLineRequest struct {
	MPPlanningHeaderID     uuid.UUID `json:"mp_planning_header_id" validate:"required"`
	OrganizationLocationID uuid.UUID `json:"organization_location_id" validate:"required"`
	JobLevelID             uuid.UUID `json:"job_level_id" validate:"required"`
	JobID                  uuid.UUID `json:"job_id" validate:"required"`
	Existing               int       `json:"existing" validate:"required"`
	Recruit                int       `json:"recruit" validate:"required"`
	SuggestedRecruit       int       `json:"suggested_recruit" validate:"required"`
	Promotion              int       `json:"promotion" validate:"required"`
	Total                  int       `json:"total" validate:"required"`
	RecruitPH              int       `json:"recruit_ph" validate:"required"`
	RecruitMT              int       `json:"recruit_mt" validate:"required"`
	// RemainingBalancePH     int       `json:"remaining_balance_ph" validate:"required"`
	// RemainingBalanceMT     int       `json:"remaining_balance_mt" validate:"required"`
}

type FindHeaderByMPPPeriodIdMPPlanningRequest struct {
	MPPPeriodID string `json:"mpp_period_id" validate:"required"`
}

type CreateOrUpdateBatchLineMPPlanningLinesRequest struct {
	MPPlanningHeaderID uuid.UUID `json:"mp_planning_header_id" validate:"required"`
	MPPlanningLines    []struct {
		ID                     uuid.UUID `json:"id" validate:"omitempty"`
		OrganizationLocationID uuid.UUID `json:"organization_location_id" validate:"required"`
		JobLevelID             uuid.UUID `json:"job_level_id" validate:"required"`
		JobID                  uuid.UUID `json:"job_id" validate:"required"`
		Existing               int       `json:"existing" validate:"required"`
		Recruit                int       `json:"recruit" validate:"required"`
		SuggestedRecruit       int       `json:"suggested_recruit" validate:"required"`
		Promotion              int       `json:"promotion" validate:"required"`
		Total                  int       `json:"total" validate:"required"`
		RecruitPH              int       `json:"recruit_ph" validate:"required"`
		RecruitMT              int       `json:"recruit_mt" validate:"required"`
		// RemainingBalancePH     int       `json:"remaining_balance_ph" validate:"required"`
		// RemainingBalanceMT     int       `json:"remaining_balance_mt" validate:"required"`
	} `json:"mp_planning_lines" validate:"required"`
}

type UpdateLineMPPlanningLineRequest struct {
	ID                     uuid.UUID `json:"id" validate:"required"`
	MPPlanningHeaderID     uuid.UUID `json:"mp_planning_header_id" validate:"required"`
	OrganizationLocationID uuid.UUID `json:"organization_location_id" validate:"required"`
	JobLevelID             uuid.UUID `json:"job_level_id" validate:"required"`
	JobID                  uuid.UUID `json:"job_id" validate:"required"`
	Existing               int       `json:"existing" validate:"required"`
	Recruit                int       `json:"recruit" validate:"required"`
	SuggestedRecruit       int       `json:"suggested_recruit" validate:"required"`
	Promotion              int       `json:"promotion" validate:"required"`
	Total                  int       `json:"total" validate:"required"`
	RecruitPH              int       `json:"recruit_ph" validate:"required"`
	RecruitMT              int       `json:"recruit_mt" validate:"required"`
	// RemainingBalancePH     int       `json:"remaining_balance_ph" validate:"required"`
	// RemainingBalanceMT     int       `json:"remaining_balance_mt" validate:"required"`
}

type DeleteLineMPPlanningLineRequest struct {
	ID string `json:"id" validate:"required"`
}
