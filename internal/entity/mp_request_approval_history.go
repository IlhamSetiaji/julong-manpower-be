package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MPRequestApprovalHistoryStatus string

const (
	MPRequestApprovalHistoryStatusApproved     MPRequestApprovalHistoryStatus = "APPROVED"
	MPRequestApprovalHistoryStatusRejected     MPRequestApprovalHistoryStatus = "REJECTED"
	MPRequestApprovalHistoryStatusNeedApproval MPRequestApprovalHistoryStatus = "NEED APPROVAL"
)

type MPRequestApprovalHistoryLevel string

const (
	MPRequestApprovalHistoryLevelStaff    MPRequestApprovalHistoryLevel = "Level Staff"
	MPRequestApprovalHistoryLevelHeadDept MPRequestApprovalHistoryLevel = "Level Head Department"
	MPRequestApprovalHistoryLevelVP       MPRequestApprovalHistoryLevel = "Level VP"
	MPRequestApprovalHistoryLevelCEO      MPRequestApprovalHistoryLevel = "Level CEO"
	MPPRequestApprovalHistoryLevelHRDHO   MPRequestApprovalHistoryLevel = "Level HRD HO"
)

type MPRequestApprovalHistory struct {
	gorm.Model        `json:"-"`
	ID                uuid.UUID                      `json:"id" gorm:"type:char(36);primaryKey;"`
	MPRequestHeaderID uuid.UUID                      `json:"mp_request_header_id" gorm:"type:char(36);"`
	ApproverID        uuid.UUID                      `json:"approver_id" gorm:"type:char(36);"`
	ApproverName      string                         `json:"approver_name" gorm:"type:varchar(255);"`
	Notes             string                         `json:"notes" gorm:"type:text;"`
	Level             string                         `json:"level" gorm:"type:varchar(255);"`
	Status            MPRequestApprovalHistoryStatus `json:"status" gorm:"not null"`

	MPRequestHeader     MPRequestHeader      `gorm:"foreignKey:MPRequestHeaderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ManpowerAttachments []ManpowerAttachment `json:"manpower_attachments" gorm:"polymorphicType:OwnerType;polymorphicId:OwnerID;polymorphicValue:mp_request_approval_histories" constraint:"OnUpdate:CASCADE,OnDelete:CASCADE"`
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
