package repo

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/chat/entity"
)

type (
	RoomRepo interface {
		Create(ctx context.Context, room *entity.Room) error
		FindByID(ctx context.Context, id string) (*entity.Room, error)
		FindByUserID(ctx context.Context, userID string) ([]*entity.Room, error)
		Patch(ctx context.Context, id string, room *entity.Room) error
		Delete(ctx context.Context, id string) error
	}
	MembershipRepo interface {
		Create(ctx context.Context, membership *entity.Membership) error
		FindByRoomID(ctx context.Context, roomID string) ([]*entity.Membership, error)
		FindByUserID(ctx context.Context, userID string) ([]*entity.Membership, error)
		FindByRoomIDAndUserID(ctx context.Context, roomID, userID string) (*entity.Membership, error)
		Delete(ctx context.Context, roomID, userID string) error
		DeleteByUserID(ctx context.Context, userID string) error
	}
)
