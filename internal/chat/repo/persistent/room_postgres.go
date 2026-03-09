package persistent

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/chat/entity"
	"gorm.io/gorm"
)

type RoomRepo struct {
	db *gorm.DB
}

func NewRoomRepo(db *gorm.DB) *RoomRepo {
	return &RoomRepo{
		db: db,
	}
}

func (r *RoomRepo) Create(ctx context.Context, room *entity.Room) error {
	db := r.db.WithContext(ctx)
	return db.Create(room).Error
}

func (r *RoomRepo) FindByID(ctx context.Context, id string) (*entity.Room, error) {
	db := r.db.WithContext(ctx)
	var room entity.Room
	if err := db.Preload("Memberships").Where("id = ?", id).First(&room).Error; err != nil {
		return &entity.Room{}, err
	}
	return &room, nil
}

func (r *RoomRepo) FindByUserID(ctx context.Context, userID string) ([]*entity.Room, error) {
	db := r.db.WithContext(ctx)
	var roomValues []entity.Room
	if err := db.
		Model(&entity.Room{}).
		Joins("JOIN memberships ON memberships.room_id = rooms.id").
		Where("memberships.user_id = ?", userID).
		Preload("Memberships", "user_id = ?", userID).
		Distinct("rooms.*").
		Find(&roomValues).Error; err != nil {
		return nil, err
	}

	rooms := make([]*entity.Room, len(roomValues))
	for i := range roomValues {
		rooms[i] = &roomValues[i]
	}
	return rooms, nil
}

func (r *RoomRepo) Patch(ctx context.Context, id string, room *entity.Room) error {
	db := r.db.WithContext(ctx)
	result := db.Model(&entity.Room{}).Where("id = ?", id).Updates(map[string]string{
		"title": room.Title,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *RoomRepo) Delete(ctx context.Context, id string) error {
	db := r.db.WithContext(ctx)
	result := db.Delete(&entity.Room{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
