package v1

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/profile/entity"
	v1 "github.com/KimNattanan/go-chat-backend/internal/profile/proto/v1"
	"github.com/google/uuid"
)

func (r *V1) CreateProfile(ctx context.Context, req *v1.CreateProfileRequest) (*v1.ProfileResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	profile := &entity.Profile{
		UserID: userID,
		Email:  req.Email,
		Name:   req.Name,
	}
	if err := r.profileUseCase.Create(ctx, profile); err != nil {
		return nil, err
	}
	return &v1.ProfileResponse{
		Profile: toProtoProfile(profile),
	}, nil
}

func (r *V1) DeleteProfile(ctx context.Context, req *v1.DeleteProfileRequest) (*v1.DeleteProfileResponse, error) {
	if err := r.profileUseCase.Delete(ctx, req.UserId); err != nil {
		return nil, err
	}
	return &v1.DeleteProfileResponse{
		Message: "profile deleted",
	}, nil
}
