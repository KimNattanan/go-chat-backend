package usecase

import (
	"context"
	"time"

	"github.com/KimNattanan/go-chat-backend/internal/auth/entity"
)

type (
	UserUseCase interface {
		Login(ctx context.Context, email, password string) (*entity.User, error)
		Register(ctx context.Context, email, password, name string) (*entity.User, error)
		FindByID(ctx context.Context, id string) (*entity.User, error)
		FindByEmail(ctx context.Context, email string) (*entity.User, error)
		Delete(ctx context.Context, userID string) error
	}
	SessionUseCase interface {
		Create(ctx context.Context, session *entity.Session) error
		FindByID(ctx context.Context, id string) (*entity.Session, error)
		FindByUserID(ctx context.Context, userID string) ([]*entity.Session, error)
		Revoke(ctx context.Context, id string) error
		Delete(ctx context.Context, id string) error
		Refresh(ctx context.Context, userID, id, newID string, expiresAt time.Time) error
	}
)
