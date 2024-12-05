package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestMajor struct {
	gorm.Model        `json:"-"`
	ID                uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	MajorID           uuid.UUID `json:"major_id" gorm:"type:char(36);"`
	MPRequestHeaderID uuid.UUID `json:"mp_request_header_id" gorm:"type:char(36);"`

	Major           Major           `json:"major" gorm:"foreignKey:MajorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MPRequestHeader MPRequestHeader `json:"mp_request_header" gorm:"foreignKey:MPRequestHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (m *RequestMajor) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(7 * time.Hour)
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *RequestMajor) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (RequestMajor) TableName() string {
	return "request_majors"
}
