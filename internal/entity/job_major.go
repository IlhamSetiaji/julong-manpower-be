package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobMajor struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	Title      string    `json:"title" gorm:"type:varchar(255);not null"`
	StartDate  time.Time `json:"start_date" gorm:"type:date;not null"`
	EndDate    time.Time `json:"end_date" gorm:"type:date;not null"`
	Status     string    `json:"status" gorm:"type:enum('open', 'close', 'complete');default:'open'"`

	JobMajors []JobMajor `json:"job_majors" gorm:"foreignKey:JobMajorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (m *JobMajor) BeforeCreate() (err error) {
	m.CreatedAt = time.Now().Add(7 * time.Hour)
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *JobMajor) BeforeUpdate() (err error) {
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *JobMajor) TableName() string {
	return "job_mayors"
}
