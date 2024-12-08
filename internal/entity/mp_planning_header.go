package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MPPlaningStatus string

const (
	MPPlaningStatusDraft    MPPlaningStatus = "DRAFT"
	MPPlaningStatusReject   MPPlaningStatus = "REJECT"
	MPPlaningStatusSubmit   MPPlaningStatus = "SUBMIT"
	MPPlaningStatusComplete MPPlaningStatus = "COMPLETE"
)

type MPPlanningHeader struct {
	gorm.Model        `json:"-"`
	ID                uuid.UUID       `json:"id" gorm:"type:char(36);primaryKey;"`
	MPPPeriodID       uuid.UUID       `json:"mpp_period_id" gorm:"type:char(36);"`
	OrganizationID    *uuid.UUID      `json:"organization_id" gorm:"type:char(36);not null;"`
	EmpOrganizationID *uuid.UUID      `json:"emp_organization_id" gorm:"type:char(36);not null;"`
	JobID             *uuid.UUID      `json:"job_id" gorm:"type:char(36);not null;"` // job_id
	DocumentNumber    string          `json:"document_number" gorm:"type:varchar(255);not null;"`
	DocumentDate      time.Time       `json:"document_date" gorm:"type:date;not null;"`
	Notes             string          `json:"notes" gorm:"type:text;default:null"`
	TotalRecruit      float64         `json:"total_recruit" gorm:"type:decimal(18,2);default:0"`
	TotalPromote      float64         `json:"total_promote" gorm:"ty	pe:decimal(18,2);default:0"`
	Status            MPPlaningStatus `json:"status" gorm:"default:'DRAFT'"`
	RecommendedBy     string          `json:"recommended_by" gorm:"type:text;"`   // free text
	ApprovedBy        string          `json:"approved_by" gorm:"type:text;"`      // free text
	RequestorID       *uuid.UUID      `json:"requestor_id" gorm:"type:char(36);"` // user_id
	NotesAttach       string          `json:"notes_attach" gorm:"type:text;"`

	MPPPeriod           MPPPeriod            `json:"mpp_period" gorm:"foreignKey:MPPPeriodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MPPlanningLines     []MPPlanningLine     `json:"mp_planning_lines" gorm:"foreignKey:MPPlanningHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ManpowerAttachments []ManpowerAttachment `json:"manpower_attachments" gorm:"polymorphicType:OwnerType;polymorphicId:OwnerID;polymorphicValue:mp_planning_headers"`
	MPRequestHeaders    []MPRequestHeader    `json:"mp_request_headers" gorm:"foreignKey:MPPlanningHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	OrganizationName    string `json:"organization_name" gorm:"-"`
	EmpOrganizationName string `json:"emp_organization_name" gorm:"-"`
	JobName             string `json:"job_name" gorm:"-"`
	RequestorName       string `json:"requestor_name" gorm:"-"`
}

func (m *MPPlanningHeader) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(time.Hour * 7)
	m.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (m *MPPlanningHeader) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (MPPlanningHeader) TableName() string {
	return "mp_planning_headers"
}
