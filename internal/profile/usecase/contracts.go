package usecase

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/profile/entity"
)

type (
	ProfileUseCase interface {
		Create(ctx context.Context, profile *entity.Profile) error
		FindByID(ctx context.Context, userID string) (*entity.Profile, error)
		Patch(ctx context.Context, userID string, profile *entity.Profile) (*entity.Profile, error)
		Delete(ctx context.Context, userID string) error
	}
)
