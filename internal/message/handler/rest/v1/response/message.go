package response

import (
	"time"

	"github.com/google/uuid"
)

type MessageResponse struct {
	ID        uuid.UUID `json:"id"`
	RoomID    uuid.UUID `json:"room_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   uuid.UUID `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
