package response

import (
	"time"

	"github.com/google/uuid"
)

type MembershipResponse struct {
	RoomID    uuid.UUID `json:"room_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
