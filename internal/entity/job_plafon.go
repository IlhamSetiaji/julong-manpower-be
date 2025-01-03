package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobPlafon struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey;"`
	JobID      *uuid.UUID `json:"job_id" gorm:"type:char(36);not null;unique"`
	Plafon     int        `json:"plafon" gorm:"type:int;default:0"`

	JobName          string `json:"job_name" gorm:"-"`
	OrganizationName string `json:"organization_name" gorm:"-"`
}

func (m *JobPlafon) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(7 * time.Hour)
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *JobPlafon) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (JobPlafon) TableName() string {
	return "job_plafons"
}
