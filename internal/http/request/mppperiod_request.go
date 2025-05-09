package request

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/entity"
	"github.com/google/uuid"
)

type FindAllPaginatedMPPPeriodRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Search   string `json:"search"`
}

type FindByIdMPPPeriodRequest struct {
	ID uuid.UUID `json:"id"`
}

type CreateMPPPeriodRequest struct {
	Title           string                 `json:"title" validate:"required"`
	StartDate       string                 `json:"start_date" validate:"required,datetime=2006-01-02,date_today_or_later"`
	EndDate         string                 `json:"end_date" validate:"required,datetime=2006-01-02,date_today_or_later"`
	BudgetStartDate string                 `json:"budget_start_date" validate:"required,datetime=2006-01-02,date_today_or_later"`
	BudgetEndDate   string                 `json:"budget_end_date" validate:"required,datetime=2006-01-02,date_today_or_later"`
	Status          entity.MPPPeriodStatus `json:"status" validate:"omitempty,MPPPeriodStatusValidation"`
}

type UpdateMPPPeriodRequest struct {
	ID              uuid.UUID              `json:"id" validate:"required"`
	Title           string                 `json:"title" validate:"required"`
	StartDate       string                 `json:"start_date" validate:"required,datetime=2006-01-02,date_today_or_later"`
	EndDate         string                 `json:"end_date" validate:"required,datetime=2006-01-02,date_today_or_later"`
	BudgetStartDate string                 `json:"budget_start_date" validate:"required,datetime=2006-01-02,date_today_or_later"`
	BudgetEndDate   string                 `json:"budget_end_date" validate:"required,datetime=2006-01-02,date_today_or_later"`
	Status          entity.MPPPeriodStatus `json:"status" validate:"omitempty,MPPPeriodStatusValidation"`
}

type FindByCurrentDateAndStatusMPPPeriodRequest struct {
	Status entity.MPPPeriodStatus `json:"status" validate:"required,MPPPeriodStatusValidation"`
}

type DeleteMPPPeriodRequest struct {
	ID uuid.UUID `json:"id"`
}
