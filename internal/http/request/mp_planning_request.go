package request

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
)

type FindAllHeadersPaginatedMPPlanningRequest struct {
	Page          int    `json:"page" binding:"required"`
	PageSize      int    `json:"pageSize" binding:"required"`
	Search        string `json:"search"`
	ApproverType  string `json:"approver_type"`
	OrgLocationID string `json:"org_location_id"`
	OrgID         string `json:"org_id"`
	Status        string `json:"status"`
	IsNull        string `json:"is_null"`
	RequestorID   string `json:"requestor_id"`
}

type MPPlanningHeaderRequest struct {
	ID                     string                 `json:"id" validate:"omitempty"`
	DocumentNumber         string                 `json:"document_number" validate:"omitempty"`
	OrganizationID         string                 `json:"organization_id" validate:"omitempty"`
	EmpOrganizationID      string                 `json:"emp_organization_id" validate:"omitempty"`
	OrganizationLocationID string                 `json:"organization_location_id" validate:"omitempty"`
	JobID                  string                 `json:"job_id" validate:"omitempty"`
	MPPPeriodID            string                 `json:"mpp_period_id" validate:"omitempty"`
	Status                 entity.MPPlaningStatus `json:"status" validate:"omitempty"`
}

type FindHeaderByIdMPPlanningRequest struct {
	ID string `json:"id" validate:"required"`
}

type ManpowerAttachmentRequest struct {
	FileName string `json:"file_name" validate:"required"`
	FilePath string `json:"file_path" validate:"required"`
	FileType string `json:"file_type" validate:"required"`
}

type UpdateStatusMPPlanningHeaderRequest struct {
	ID          string                                `json:"id" validate:"required"`
	Status      entity.MPPlaningStatus                `json:"status" validate:"required,MPPlaningStatusValidation"`
	Notes       string                                `json:"notes" validate:"omitempty"`
	Level       entity.MPPlanningApprovalHistoryLevel `json:"level" validate:"required,MPPlanningApprovalHistoryLevelValidation"`
	Attachments []ManpowerAttachmentRequest           `json:"attachments" validate:"omitempty,dive"`
	ApprovedBy  string                                `json:"approved_by" validate:"required"`
	ApproverID  uuid.UUID                             `json:"approver_id" validate:"required"`
	// ApproverName string                      `json:"approved_by_name" validate:"omitempty"`
}

type UpdateStatusPartialMPPlanningHeaderRequest struct {
	ApproverID uuid.UUID                    `json:"approver_id" validate:"required"`
	Payload    []UpdateStatusPartialPayload `json:"payload" validate:"required,dive"`
}

type UpdateStatusPartialPayload struct {
	ID    string `json:"id" validate:"required,uuid"`
	Notes string `json:"notes" validate:"omitempty"`
}

type CreateHeaderMPPlanningRequest struct {
	MPPPeriodID            uuid.UUID                   `json:"mpp_period_id" validate:"required"`
	OrganizationID         uuid.UUID                   `json:"organization_id" validate:"required"`
	EmpOrganizationID      uuid.UUID                   `json:"emp_organization_id" validate:"required"`
	OrganizationLocationID uuid.UUID                   `json:"organization_location_id" validate:"required"` // organization_location_id
	JobID                  uuid.UUID                   `json:"job_id" validate:"required"`                   // job_id
	DocumentNumber         string                      `json:"document_number" validate:"required"`
	DocumentDate           string                      `json:"document_date" validate:"required,datetime=2006-01-02,date_today_or_later"`
	Notes                  string                      `json:"notes" validate:"omitempty"`
	TotalRecruit           float64                     `json:"total_recruit" validate:"omitempty"`
	TotalPromote           float64                     `json:"total_promote" validate:"omitempty"`
	Status                 entity.MPPlaningStatus      `json:"status" validate:"required,MPPlaningStatusValidation"`
	RecommendedBy          string                      `json:"recommended_by" validate:"omitempty"`
	ApprovedBy             string                      `json:"approved_by" validate:"omitempty"`
	RequestorID            uuid.UUID                   `json:"requestor_id" validate:"required"`
	NotesAttach            string                      `json:"notes_attach" validate:"omitempty"`
	Attachments            []ManpowerAttachmentRequest `json:"attachments" validate:"omitempty,dive"`
}

type UpdateHeaderMPPlanningRequest struct {
	ID                     uuid.UUID                   `json:"id" validate:"required"`
	MPPPeriodID            uuid.UUID                   `json:"mpp_period_id" validate:"required"`
	OrganizationID         uuid.UUID                   `json:"organization_id" validate:"required"`
	EmpOrganizationID      uuid.UUID                   `json:"emp_organization_id" validate:"required"`
	OrganizationLocationID uuid.UUID                   `json:"organization_location_id" validate:"required"`
	JobID                  uuid.UUID                   `json:"job_id" validate:"required"` // job_id
	DocumentNumber         string                      `json:"document_number" validate:"required"`
	DocumentDate           string                      `json:"document_date" validate:"required,datetime=2006-01-02"`
	Notes                  string                      `json:"notes" validate:"omitempty"`
	TotalRecruit           float64                     `json:"total_recruit" validate:"omitempty"`
	TotalPromote           float64                     `json:"total_promote" validate:"omitempty"`
	Status                 entity.MPPlaningStatus      `json:"status" validate:"required,MPPlaningStatusValidation"`
	RecommendedBy          string                      `json:"recommended_by" validate:"omitempty"`
	ApprovedBy             string                      `json:"approved_by" validate:"omitempty"`
	RequestorID            uuid.UUID                   `json:"requestor_id" validate:"required"`
	NotesAttach            string                      `json:"notes_attach" validate:"omitempty"`
	Attachments            []ManpowerAttachmentRequest `json:"attachments" validate:"omitempty,dive"`
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
	OrganizationLocationID uuid.UUID `json:"organization_location_id" validate:"omitempty"` // organization_location_id
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
		OrganizationLocationID uuid.UUID `json:"organization_location_id" validate:"omitempty"` // organization_location_id
		JobLevelID             uuid.UUID `json:"job_level_id" validate:"required"`
		JobID                  uuid.UUID `json:"job_id" validate:"required"`
		Existing               int       `json:"existing" validate:"required"`
		Recruit                int       `json:"recruit" validate:"required"`
		SuggestedRecruit       int       `json:"suggested_recruit" validate:"required"`
		Promotion              int       `json:"promotion" validate:"required"`
		Total                  int       `json:"total" validate:"required"`
		RecruitPH              int       `json:"recruit_ph" validate:"required"`
		RecruitMT              int       `json:"recruit_mt" validate:"required"`
		IsCreate               bool      `json:"is_create" validate:"omitempty"`
		// RemainingBalancePH     int       `json:"remaining_balance_ph" validate:"required"`
		// RemainingBalanceMT     int       `json:"remaining_balance_mt" validate:"required"`
	} `json:"mp_planning_lines" validate:"required"`
	DeletedLineIDs []string `json:"deleted_line_ids" validate:"omitempty,dive"`
}

type UpdateLineMPPlanningLineRequest struct {
	ID                     uuid.UUID `json:"id" validate:"required"`
	MPPlanningHeaderID     uuid.UUID `json:"mp_planning_header_id" validate:"required"`
	OrganizationLocationID uuid.UUID `json:"organization_location_id" validate:"omitempty"` // organization_location_id
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
