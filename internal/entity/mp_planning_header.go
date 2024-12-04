package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MPPlaningStatus string

const (
	MPPlaningStatusOpen     MPPlaningStatus = "OPEN"
	MPPlaningStatusClose    MPPlaningStatus = "CLOSE"
	MPPlaningStatusComplete MPPlaningStatus = "COMPLETE"
)

type MPPlanningHeader struct {
	gorm.Model        `json:"-"`
	ID                uuid.UUID       `json:"id" gorm:"type:char(36);primaryKey;"`
	MPPPeriodID       uuid.UUID       `json:"mpp_period_id" gorm:"type:char(36);"`
	OrganizationID    *uuid.UUID      `json:"organization_id" gorm:"type:char(36);not null;"`
	EmpOrganizationID *uuid.UUID      `json:"emp_organization_id" gorm:"type:char(36);not null;"`
	DocumentNumber    string          `json:"document_number" gorm:"type:varchar(255);not null;"`
	DocumentDate      time.Time       `json:"document_date" gorm:"type:date;not null;"`
	Notes             string          `json:"notes" gorm:"type:text;"`
	TotalRecruit      float64         `json:"total_recruit" gorm:"type:decimal(18,2);default:0"`
	TotalPromote      float64         `json:"total_promote" gorm:"type:decimal(18,2);default:0"`
	Status            MPPlaningStatus `json:"status" gorm:"default:'open'"`
	RecommendedBy     *uuid.UUID      `json:"recommended_by" gorm:"type:char(36);"`
	ApprovedBy        *uuid.UUID      `json:"approved_by" gorm:"type:char(36);"`
	RequestorID       *uuid.UUID      `json:"requestor_id" gorm:"type:char(36);"`

	MPPPeriod       MPPPeriod        `json:"mpp_period" gorm:"foreignKey:MPPPeriodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MPPlanningLines []MPPlanningLine `json:"mp_planning_lines" gorm:"foreignKey:MPPlanningHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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
