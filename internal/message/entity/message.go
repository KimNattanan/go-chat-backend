package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RoomID    uuid.UUID `gorm:"type:uuid" json:"room_id"`
	UserID    uuid.UUID `gorm:"type:uuid" json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Message) BeforeCreate(db *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}
