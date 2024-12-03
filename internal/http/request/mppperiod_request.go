package request

import "github.com/google/uuid"

type FindAllPaginatedMPPPeriodRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Search   string `json:"search"`
}

type FindByIdMPPPeriodRequest struct {
	ID uuid.UUID `json:"id"`
}
