package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestType struct {
	gorm.Model
	ID            uuid.UUID `json:"id" gorm:"type:char(32);primaryKey;"`
	Name          string    `json:"name" gorm:"type:varchar(255);not null"`
	IsReplacement bool      `json:"is_replacement" gorm:"type:boolean;default:false"`

	MPRequestHeaders []MPRequestHeader `json:"mp_request_headers" gorm:"foreignKey:RequestTypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (m *RequestType) BeforeCreate() (err error) {
	m.ID = uuid.New()
	m.CreatedAt = m.CreatedAt.Add(7 * 60 * 60)
	m.UpdatedAt = m.UpdatedAt.Add(7 * 60 * 60)
	return nil
}

func (m *RequestType) BeforeUpdate() (err error) {
	m.UpdatedAt = m.UpdatedAt.Add(7 * 60 * 60)
	return nil
}

func (RequestType) TableName() string {
	return "request_types"
}
