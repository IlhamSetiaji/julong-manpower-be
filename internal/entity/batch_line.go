package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BatchLine struct {
	gorm.Model             `json:"-"`
	ID                     uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	BatchHeaderID          uuid.UUID `json:"batch_header_id" gorm:"type:char(36);not null;"`
	MPPlanningHeaderID     uuid.UUID `json:"mp_planning_header_id" gorm:"type:char(36);not null;"`
	OrganizationID         uuid.UUID `json:"organization_id" gorm:"type:char(36);not null;"`
	OrganizationLocationID uuid.UUID `json:"organization_location_id" gorm:"type:char(36);not null;"`

	OrganizationName         string `json:"organization_name" gorm:"-"`
	OrganizationLocationName string `json:"organization_location_name" gorm:"-"`

	BatchHeader      BatchHeader      `json:"batch_header" gorm:"foreignKey:BatchHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MPPlanningHeader MPPlanningHeader `json:"mp_planning_header" gorm:"foreignKey:MPPlanningHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (m *BatchLine) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(time.Hour * 7)
	m.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (m *BatchLine) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (BatchLine) TableName() string {
	return "batch_lines"
}
