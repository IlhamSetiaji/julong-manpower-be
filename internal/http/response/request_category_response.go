package response

import (
	"time"

	"github.com/google/uuid"
)

type RequestCategoryResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	IsReplacement bool      `json:"is_replacement"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
