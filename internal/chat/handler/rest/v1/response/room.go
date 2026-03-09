package response

import (
	"time"

	"github.com/google/uuid"
)

type RoomResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Memberships []*struct {
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"memberships"`
}
