package v1

import (
	"context"

	v1 "github.com/KimNattanan/go-chat-backend/internal/auth/proto/v1"
)

func (r *V1) FindUserByID(ctx context.Context, req *v1.FindUserByIDRequest) (*v1.UserResponse, error) {
	return nil, nil
}

func (r *V1) FindUserByEmail(ctx context.Context, req *v1.FindUserByEmailRequest) (*v1.UserResponse, error) {
	return nil, nil
}

func (r *V1) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {
	return nil, nil
}
