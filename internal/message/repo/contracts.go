package repo

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/message/entity"
)

type (
	MessageRepo interface {
		Create(ctx context.Context, message *entity.Message) error
		FindByID(ctx context.Context, id string) (*entity.Message, error)
		FindByRoomID(ctx context.Context, roomID string) ([]*entity.Message, error)
		FindByUserID(ctx context.Context, userID string) ([]*entity.Message, error)
		FindByRoomIDAndUserID(ctx context.Context, roomID, userID string) ([]*entity.Message, error)
		AnonymizeUserMessages(ctx context.Context, userID string) error
		Delete(ctx context.Context, id string) error
		DeleteByRoomID(ctx context.Context, roomID string) error
	}
)
