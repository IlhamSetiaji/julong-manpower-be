package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MPPlanningLine struct {
	gorm.Model             `json:"-"`
	ID                     uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey;"`
	MPPlanningHeaderID     uuid.UUID  `json:"mp_planning_header_id" gorm:"type:char(36);"`
	OrganizationLocationID *uuid.UUID `json:"organization_location_id" gorm:"type:char(36);not null"`
	JobLevelID             *uuid.UUID `json:"job_level_id" gorm:"type:char(36);not null"`
	JobID                  *uuid.UUID `json:"job_id" gorm:"type:char(36);not null"`
	Existing               int        `json:"existing" gorm:"type:int;default:0"`
	Recruit                int        `json:"recruit" gorm:"type:int;default:0"`
	SuggestedRecruit       int        `json:"suggested_recruit" gorm:"type:int;default:0"`
	Promotion              int        `json:"promotion" gorm:"type:int;default:0"`
	Total                  int        `json:"total" gorm:"type:int;default:0"`
	RemainingBalance       int        `json:"remaining_balance" gorm:"type:int;default:0"`

	MPPlanningHeader MPPlanningHeader `json:"mp_planning_header" gorm:"foreignKey:MPPlanningHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (m *MPPlanningLine) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(7 * time.Hour)
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *MPPlanningLine) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (MPPlanningLine) TableName() string {
	return "mp_planning_lines"
}
