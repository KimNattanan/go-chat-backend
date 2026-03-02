package entity

import (
	"time"

	"github.com/google/uuid"
)

type Membership struct {
	RoomID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"room_id"`
	UserID    uuid.UUID `gorm:"type:uuid;primarykey" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
