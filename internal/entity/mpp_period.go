package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MPPPeriodStatus string

const (
	MPPeriodStatusOpen     MPPPeriodStatus = "open"
	MPPeriodStatusClose    MPPPeriodStatus = "close"
	MPPeriodStatusComplete MPPPeriodStatus = "complete"
)

type MPPPeriod struct {
	gorm.Model
	ID        uuid.UUID       `json:"id" gorm:"type:char(32);primaryKey;"`
	Title     string          `json:"title"`
	StartDate time.Time       `json:"start_date" gorm:"type:date"`
	EndDate   time.Time       `json:"end_date" gorm:"type:date"`
	Status    MPPPeriodStatus `json:"status" gorm:"default:'open'"`

	MPPlanningHeaders []MPPlanningHeader `json:"mp_planning_headers" gorm:"foreignKey:MPPPeriodID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (m *MPPPeriod) BeforeCreate() (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(7 * time.Hour)
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *MPPPeriod) BeforeUpdate() (err error) {
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *MPPPeriod) TableName() string {
	return "mpp_periods"
}
