package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestCategory struct {
	gorm.Model    `json:"-"`
	ID            uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;"`
	Name          string    `json:"name" gorm:"type:varchar(255);not null;"`
	IsReplacement bool      `json:"is_replacement" gorm:"type:boolean;default:false;"`

	MPRequestHeaders []MPRequestHeader `json:"mp_request_headers" gorm:"foreignKey:RequestCategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (m *RequestCategory) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(time.Hour * 7)
	m.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (m *RequestCategory) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(time.Hour * 7)
	return nil
}

func (RequestCategory) TableName() string {
	return "request_categories"
}
