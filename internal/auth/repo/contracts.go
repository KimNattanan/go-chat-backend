package repo

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/auth/entity"
)

type (
	UserUseCase interface {
		Create(ctx context.Context, user *entity.User) error
		FindByID(ctx context.Context, id string) (*entity.User, error)
		FindByEmail(ctx context.Context, email string) (*entity.User, error)
		Delete(ctx context.Context, userID string) error
	}
	SessionRepo interface {
		Create(ctx context.Context, session *entity.Session) error
		FindByID(ctx context.Context, id string) (*entity.Session, error)
		FindByUserID(ctx context.Context, userID string) ([]*entity.Session, error)
		Revoke(ctx context.Context, id string) error
		Delete(ctx context.Context, id string) error
	}
)
