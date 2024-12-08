package response

import (
	"time"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
)

type MPPlanningHeaderResponse struct {
	ID                uuid.UUID              `json:"id"`
	MPPPeriodID       uuid.UUID              `json:"mpp_period_id"`
	OrganizationID    *uuid.UUID             `json:"organization_id"`
	EmpOrganizationID *uuid.UUID             `json:"emp_organization_id"`
	JobID             *uuid.UUID             `json:"job_id"` // job_id
	DocumentNumber    string                 `json:"document_number"`
	DocumentDate      time.Time              `json:"document_date"`
	Notes             string                 `json:"notes"`
	TotalRecruit      float64                `json:"total_recruit"`
	TotalPromote      float64                `json:"total_promote"`
	Status            entity.MPPlaningStatus `json:"status"`
	RecommendedBy     string                 `json:"recommended_by"` // free text
	ApprovedBy        string                 `json:"approved_by"`    // free text
	RequestorID       *uuid.UUID             `json:"requestor_id"`   // user_id
	NotesAttach       string                 `json:"notes_attach"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         *time.Time             `json:"deleted_at"`

	OrganizationName    string                    `json:"organization_name" gorm:"-"`
	EmpOrganizationName string                    `json:"emp_organization_name" gorm:"-"`
	JobName             string                    `json:"job_name" gorm:"-"`
	RequestorName       string                    `json:"requestor_name" gorm:"-"`
	MPPPeriod           *MPPeriodResponse         `json:"mpp_period"`
	MPPlanningLines     []*MPPlanningLineResponse `json:"mp_planning_lines"`
}

type FindHeaderByMPPPeriodIdMPPlanningResponse struct {
	ID                uuid.UUID              `json:"id"`
	MPPPeriodID       uuid.UUID              `json:"mpp_period_id"`
	OrganizationID    *uuid.UUID             `json:"organization_id"`
	EmpOrganizationID *uuid.UUID             `json:"emp_organization_id"`
	JobID             *uuid.UUID             `json:"job_id"` // job_id
	DocumentNumber    string                 `json:"document_number"`
	DocumentDate      time.Time              `json:"document_date"`
	Notes             string                 `json:"notes"`
	TotalRecruit      float64                `json:"total_recruit"`
	TotalPromote      float64                `json:"total_promote"`
	Status            entity.MPPlaningStatus `json:"status"`
	RecommendedBy     string                 `json:"recommended_by"` // free text
	ApprovedBy        string                 `json:"approved_by"`    // free text
	RequestorID       *uuid.UUID             `json:"requestor_id"`   // user_id
	NotesAttach       string                 `json:"notes_attach"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         *time.Time             `json:"deleted_at"`

	OrganizationName    string                    `json:"organization_name" gorm:"-"`
	EmpOrganizationName string                    `json:"emp_organization_name" gorm:"-"`
	JobName             string                    `json:"job_name" gorm:"-"`
	RequestorName       string                    `json:"requestor_name" gorm:"-"`
	MPPPeriod           *MPPeriodResponse         `json:"mpp_period"`
	MPPlanningLines     []*MPPlanningLineResponse `json:"mp_planning_lines"`
}

type MPPlanningLineResponse struct {
	ID                     uuid.UUID `json:"id"`
	MPPlanningHeaderID     uuid.UUID `json:"mp_planning_header_id"`
	OrganizationLocationID uuid.UUID `json:"organization_location_id"`
	JobLevelID             uuid.UUID `json:"job_level_id"`
	JobID                  uuid.UUID `json:"job_id"`
	Existing               int       `json:"existing"`
	Recruit                int       `json:"recruit"`
	SuggestedRecruit       int       `json:"suggested_recruit"`
	Promotion              int       `json:"promotion"`
	Total                  int       `json:"total"`
	RemainingBalancePH     int       `json:"remaining_balance_ph"`
	RemainingBalanceMT     int       `json:"remaining_balance_mt"`
	RecruitPH              int       `json:"recruit_ph"`
	RecruitMT              int       `json:"recruit_mt"`

	OrganizationLocationName string `json:"organization_location_name"`
	JobLevelName             string `json:"job_level_name"`
	JobName                  string `json:"job_name"`
}

type FindAllHeadersPaginatedMPPlanningResponse struct {
	MPPlanningHeaders []*MPPlanningHeaderResponse `json:"mp_planning_headers"`
	Total             int64                       `json:"total"`
}

type FindByIdMPPlanningResponse struct {
	// MPPlanningHeader *entity.MPPlanningHeader `json:"mp_planning_header"`
	// MPPlanningHeader *MPPlanningHeaderResponse `json:"mp_planning_header"`
	ID                uuid.UUID              `json:"id"`
	MPPPeriodID       uuid.UUID              `json:"mpp_period_id"`
	OrganizationID    *uuid.UUID             `json:"organization_id"`
	EmpOrganizationID *uuid.UUID             `json:"emp_organization_id"`
	JobID             *uuid.UUID             `json:"job_id"` // job_id
	DocumentNumber    string                 `json:"document_number"`
	DocumentDate      time.Time              `json:"document_date"`
	Notes             string                 `json:"notes"`
	TotalRecruit      float64                `json:"total_recruit"`
	TotalPromote      float64                `json:"total_promote"`
	Status            entity.MPPlaningStatus `json:"status"`
	RecommendedBy     string                 `json:"recommended_by"` // free text
	ApprovedBy        string                 `json:"approved_by"`    // free text
	RequestorID       *uuid.UUID             `json:"requestor_id"`   // user_id
	NotesAttach       string                 `json:"notes_attach"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         *time.Time             `json:"deleted_at"`

	OrganizationName    string                    `json:"organization_name" gorm:"-"`
	EmpOrganizationName string                    `json:"emp_organization_name" gorm:"-"`
	JobName             string                    `json:"job_name" gorm:"-"`
	RequestorName       string                    `json:"requestor_name" gorm:"-"`
	MPPPeriod           *MPPeriodResponse         `json:"mpp_period"`
	MPPlanningLines     []*MPPlanningLineResponse `json:"mp_planning_lines"`
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
	RemainingBalancePH     int       `json:"remaining_balance_ph"`
	RemainingBalanceMT     int       `json:"remaining_balance_mt"`
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
	RemainingBalanceMT     int       `json:"remaining_balance_mt"`
	RemainingBalancePH     int       `json:"remaining_balance_ph"`
	RecruitPH              int       `json:"recruit_ph"`
	RecruitMT              int       `json:"recruit_mt"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	DeletedAt              time.Time `json:"deleted_at"`
}
