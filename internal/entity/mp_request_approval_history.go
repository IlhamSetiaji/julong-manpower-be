package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MPRequestApprovalHistoryStatus string

const (
	MPRequestApprovalHistoryStatusApproved MPRequestApprovalHistoryStatus = "Approved"
	MPRequestApprovalHistoryStatusRejected MPRequestApprovalHistoryStatus = "Rejected"
)

type MPRequestApprovalHistory struct {
	gorm.Model        `json:"-"`
	ID                uuid.UUID                      `json:"id" gorm:"type:char(36);primaryKey;"`
	MPRequestHeaderID uuid.UUID                      `json:"mp_request_header_id" gorm:"type:char(36);"`
	ApproverID        uuid.UUID                      `json:"approver_id" gorm:"type:char(36);"`
	ApproverName      string                         `json:"approver_name" gorm:"type:varchar(255);"`
	Status            MPRequestApprovalHistoryStatus `json:"status" gorm:"not null"`
}

func (m *MPRequestApprovalHistory) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	m.CreatedAt = time.Now().Add(7 * time.Hour)
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (m *MPRequestApprovalHistory) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now().Add(7 * time.Hour)
	return nil
}

func (MPRequestApprovalHistory) TableName() string {
	return "mp_request_approval_histories"
}
