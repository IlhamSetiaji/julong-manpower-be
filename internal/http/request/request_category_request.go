package request

import "github.com/google/uuid"

type FindByIdRequestCategoryRequest struct {
	ID uuid.UUID `json:"id" binding:"required"`
}
