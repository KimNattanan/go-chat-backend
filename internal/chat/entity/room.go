package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Memberships []Membership `gorm:"foreignKey:RoomID"`
}

func (r *Room) BeforeCreate(db *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
