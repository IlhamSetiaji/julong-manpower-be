package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Major struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	Major      string    `json:"major" gorm:"type:varchar(255);not null;"`

	RequestMajors []RequestMajor `json:"request_majors" gorm:"foreignKey:MajorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (m *Major) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(7 * time.Hour)
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *Major) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (Major) TableName() string {
	return "majors"
}
