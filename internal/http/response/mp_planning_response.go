package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
)

type FindAllHeadersPaginatedMPPlanningResponse struct {
	MPPlanningHeaders *[]entity.MPPlanningHeader `json:"mp_planning_headers"`
	Total             int64                      `json:"total"`
}

type FindByIdMPPlanningResponse struct {
	MPPlanningHeader *entity.MPPlanningHeader `json:"mp_planning_header"`
}

type CreateMPPlanningResponse struct {
	ID                string                 `json:"id"`
	MPPPeriodID       string                 `json:"mpp_period_id"`
	OrganizationID    string                 `json:"organization_id"`
	EmpOrganizationID string                 `json:"emp_organization_id"`
	JobID             string                 `json:"job_id"`
	DocumentNumber    string                 `json:"document_number"`
	DocumentDate      time.Time              `json:"document_date"`
	Notes             string                 `json:"notes"`
	TotalRecruit      float64                `json:"total_recruit"`
	TotalPromote      float64                `json:"total_promote"`
	Status            entity.MPPlaningStatus `json:"status"`
	RecommendedBy     string                 `json:"recommended_by"`
	ApprovedBy        string                 `json:"approved_by"`
	RequestorID       string                 `json:"requestor_id"`
	NotesAttach       string                 `json:"notes_attach"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         time.Time              `json:"deleted_at"`
}

type UpdateMPPlanningResponse struct {
	ID                string                 `json:"id"`
	MPPPeriodID       string                 `json:"mpp_period_id"`
	OrganizationID    string                 `json:"organization_id"`
	EmpOrganizationID string                 `json:"emp_organization_id"`
	JobID             string                 `json:"job_id"`
	DocumentNumber    string                 `json:"document_number"`
	DocumentDate      time.Time              `json:"document_date"`
	Notes             string                 `json:"notes"`
	TotalRecruit      float64                `json:"total_recruit"`
	TotalPromote      float64                `json:"total_promote"`
	Status            entity.MPPlaningStatus `json:"status"`
	RecommendedBy     string                 `json:"recommended_by"`
	ApprovedBy        string                 `json:"approved_by"`
	RequestorID       string                 `json:"requestor_id"`
	NotesAttach       string                 `json:"notes_attach"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         time.Time              `json:"deleted_at"`
}

type FindAllLinesByHeaderIdPaginatedMPPlanningLineResponse struct {
	MPPlanningLines *[]entity.MPPlanningLine `json:"mp_planning_lines"`
	Total           int64                    `json:"total"`
}

type FindByIdMPPlanningLineResponse struct {
	MPPlanningLine *entity.MPPlanningLine `json:"mp_planning_line"`
}

type CreateMPPlanningLineResponse struct {
	ID                     string    `json:"id"`
	MPPlanningHeaderID     string    `json:"mp_planning_header_id"`
	OrganizationLocationID string    `json:"organization_location_id"`
	JobLevelID             string    `json:"job_level_id"`
	JobID                  string    `json:"job_id"`
	Existing               int       `json:"existing"`
	Recruit                int       `json:"recruit"`
	SuggestedRecruit       int       `json:"suggested_recruit"`
	Promotion              int       `json:"promotion"`
	Total                  int       `json:"total"`
	RemainingBalance       int       `json:"remaining_balance"`
	RecruitPH              int       `json:"recruit_ph"`
	RecruitMT              int       `json:"recruit_mt"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	DeletedAt              time.Time `json:"deleted_at"`
}

type UpdateMPPlanningLineResponse struct {
	ID                     string    `json:"id"`
	MPPlanningHeaderID     string    `json:"mp_planning_header_id"`
	OrganizationLocationID string    `json:"organization_location_id"`
	JobLevelID             string    `json:"job_level_id"`
	JobID                  string    `json:"job_id"`
	Existing               int       `json:"existing"`
	Recruit                int       `json:"recruit"`
	SuggestedRecruit       int       `json:"suggested_recruit"`
	Promotion              int       `json:"promotion"`
	Total                  int       `json:"total"`
	RemainingBalance       int       `json:"remaining_balance"`
	RecruitPH              int       `json:"recruit_ph"`
	RecruitMT              int       `json:"recruit_mt"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	DeletedAt              time.Time `json:"deleted_at"`
}
