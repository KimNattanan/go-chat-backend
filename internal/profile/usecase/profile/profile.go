package profile

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/profile/entity"
	"github.com/KimNattanan/go-chat-backend/internal/profile/repo"
)

type UseCase struct {
	profileRepo repo.ProfileRepo
}

func New(profileRepo repo.ProfileRepo) *UseCase {
	return &UseCase{
		profileRepo: profileRepo,
	}
}

func (u *UseCase) Create(ctx context.Context, profile *entity.Profile) error {
	return u.profileRepo.Create(ctx, profile)
}

func (u *UseCase) FindByID(ctx context.Context, userID string) (*entity.Profile, error) {
	return u.profileRepo.FindByID(ctx, userID)
}

func (u *UseCase) Patch(ctx context.Context, userID string, profile *entity.Profile) (*entity.Profile, error) {
	if err := u.profileRepo.Patch(ctx, userID, profile); err != nil {
		return nil, err
	}
	return u.profileRepo.FindByID(ctx, userID)
}

func (u *UseCase) Delete(ctx context.Context, userID string) error {
	return u.profileRepo.Delete(ctx, userID)
}
