package request

import "github.com/google/uuid"

type FindByIdMajorRequest struct {
	ID uuid.UUID `json:"id" validate:"required,uuid"`
}
