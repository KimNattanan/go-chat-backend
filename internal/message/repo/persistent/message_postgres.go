package persistent

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/message/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepo struct {
	db *gorm.DB
}

func NewMessageRepo(db *gorm.DB) *MessageRepo {
	return &MessageRepo{
		db: db,
	}
}

func (r *MessageRepo) Create(ctx context.Context, message *entity.Message) error {
	db := r.db.WithContext(ctx)
	return db.Create(message).Error
}

func (r *MessageRepo) FindByID(ctx context.Context, id string) (*entity.Message, error) {
	db := r.db.WithContext(ctx)
	var message entity.Message
	if err := db.First(&message, "id = ?", id).Error; err != nil {
		return &entity.Message{}, err
	}
	return &message, nil
}

func (r *MessageRepo) FindByRoomID(ctx context.Context, roomID string) ([]*entity.Message, error) {
	db := r.db.WithContext(ctx)
	var messageValues []entity.Message
	if err := db.Find(&messageValues, "room_id = ?", roomID).Error; err != nil {
		return nil, err
	}

	messages := make([]*entity.Message, len(messageValues))
	for i := range messageValues {
		messages[i] = &messageValues[i]
	}
	return messages, nil
}

func (r *MessageRepo) FindByUserID(ctx context.Context, userID string) ([]*entity.Message, error) {
	db := r.db.WithContext(ctx)
	var messageValues []entity.Message
	if err := db.Find(&messageValues, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	messages := make([]*entity.Message, len(messageValues))
	for i := range messageValues {
		messages[i] = &messageValues[i]
	}
	return messages, nil
}

func (r *MessageRepo) AnonymizeUserMessages(ctx context.Context, userID string) error {
	db := r.db.WithContext(ctx)
	return db.Model(&entity.Message{}).Where("user_id = ?", userID).Update("user_id", uuid.Nil).Error
}

func (r *MessageRepo) Delete(ctx context.Context, id string) error {
	db := r.db.WithContext(ctx)
	result := db.Delete(&entity.Message{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *MessageRepo) DeleteByRoomID(ctx context.Context, roomID string) error {
	db := r.db.WithContext(ctx)
	return db.Delete(&entity.Message{}, "room_id = ?", roomID).Error
}
