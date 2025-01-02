package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MPPlaningStatus string

const (
	MPPlaningStatusDraft        MPPlaningStatus = "DRAFTED"
	MPPlaningStatusReject       MPPlaningStatus = "REJECTED"
	MPPlaningStatusSubmit       MPPlaningStatus = "SUBMITTED"
	MPPlaningStatusApproved     MPPlaningStatus = "APPROVED"
	MPPlaningStatusComplete     MPPlaningStatus = "COMPLETED"
	MPPlanningStatusInProgress  MPPlaningStatus = "IN_PROGRESS"
	MPPlaningStatusNeedApproval MPPlaningStatus = "NEED APPROVAL"
)

type MPPlanningHeader struct {
	gorm.Model             `json:"-"`
	ID                     uuid.UUID       `json:"id" gorm:"type:char(36);primaryKey;"`
	MPPPeriodID            uuid.UUID       `json:"mpp_period_id" gorm:"type:char(36); not null;"`
	OrganizationID         *uuid.UUID      `json:"organization_id" gorm:"type:char(36);not null;"`
	EmpOrganizationID      *uuid.UUID      `json:"emp_organization_id" gorm:"type:char(36);not null;"`
	JobID                  *uuid.UUID      `json:"job_id" gorm:"type:char(36);not null; not null;"` // job_id
	OrganizationLocationID *uuid.UUID      `json:"organization_location_id" gorm:"type:char(36);not null;"`
	DocumentNumber         string          `json:"document_number" gorm:"type:varchar(255);not null;unique"`
	DocumentDate           time.Time       `json:"document_date" gorm:"type:date;not null;"`
	Notes                  string          `json:"notes" gorm:"type:text;default:null"`
	TotalRecruit           float64         `json:"total_recruit" gorm:"type:decimal(18,2);default:0"`
	TotalPromote           float64         `json:"total_promote" gorm:"ty	pe:decimal(18,2);default:0"`
	Status                 MPPlaningStatus `json:"status" gorm:"default:'DRAFT'"`
	RecommendedBy          string          `json:"recommended_by" gorm:"type:text;"`          // free text
	ApprovedBy             string          `json:"approved_by" gorm:"type:text;default:null"` // free text
	RequestorID            *uuid.UUID      `json:"requestor_id" gorm:"type:char(36);"`        // user_id
	NotesAttach            string          `json:"notes_attach" gorm:"type:text;"`
	ApproverManagerID      *uuid.UUID      `json:"approver_manager_id" gorm:"type:char(36);"` // user_id
	NotesManager           string          `json:"notes_manager" gorm:"type:text;"`
	ApproverRecruitmentID  *uuid.UUID      `json:"approver_recruitment_id" gorm:"type:char(36);"` // user_id
	NotesRecruitment       string          `json:"notes_recruitment" gorm:"type:text;"`
	CreatedAt              time.Time       `json:"created_at" gorm:"autoCreateTime"`

	MPPPeriod                   MPPPeriod                   `json:"mpp_period" gorm:"foreignKey:MPPPeriodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MPPlanningLines             []MPPlanningLine            `json:"mp_planning_lines" gorm:"foreignKey:MPPlanningHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ManpowerAttachments         []ManpowerAttachment        `json:"manpower_attachments" gorm:"polymorphicType:OwnerType;polymorphicId:OwnerID;polymorphicValue:mp_planning_headers" constraint:"OnUpdate:CASCADE,OnDelete:CASCADE"`
	MPRequestHeaders            []MPRequestHeader           `json:"mp_request_headers" gorm:"foreignKey:MPPlanningHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MPPlanningApprovalHistories []MPPlanningApprovalHistory `json:"mp_planning_approval_histories" gorm:"foreignKey:MPPlanningHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BatchLines                  []BatchLine                 `json:"batch_lines" gorm:"foreignKey:MPPlanningHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	OrganizationName         string `json:"organization_name" gorm:"-"`
	EmpOrganizationName      string `json:"emp_organization_name" gorm:"-"`
	JobName                  string `json:"job_name" gorm:"-"`
	RequestorName            string `json:"requestor_name" gorm:"-"`
	OrganizationLocationName string `json:"organization_location_name" gorm:"-"`
	ApproverCEOName          string `json:"approver_ceo_name" gorm:"-"`
	ApproverManagerName      string `json:"approver_manager_name" gorm:"-"`
	ApproverRecruitmentName  string `json:"approver_recruitment_name" gorm:"-"`
}

func (m *MPPlanningHeader) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	// m.CreatedAt = time.Now().Add(time.Hour * 7)
	// m.UpdatedAt = time.Now().Add(time.Hour * 7)
	m.UpdatedAt = time.Now()
	m.CreatedAt = m.UpdatedAt
	return nil
}

func (m *MPPlanningHeader) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (m *MPPlanningHeader) BeforeDelete(tx *gorm.DB) (err error) {
	if m.DeletedAt.Valid {
		return nil
	}

	randomString := uuid.New().String()

	m.DocumentNumber = m.DocumentNumber + "_deleted" + randomString

	if err := tx.Model(&m).Where("id = ?", m.ID).Updates((map[string]interface{}{
		"document_number": m.DocumentNumber,
	})).Error; err != nil {
		return err
	}

	return nil
}

func (MPPlanningHeader) TableName() string {
	return "mp_planning_headers"
}
