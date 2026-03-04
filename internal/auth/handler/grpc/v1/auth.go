package v1

import (
	"context"

	v1 "github.com/KimNattanan/go-chat-backend/internal/auth/proto/v1"
	"github.com/KimNattanan/go-chat-backend/pkg/apperror"
)

func (r *V1) FindUserByID(ctx context.Context, req *v1.FindUserByIDRequest) (*v1.UserResponse, error) {
	user, err := r.authUseCase.FindUserByID(ctx, req.Id)
	if err != nil {
		r.l.Error(err, "grpc - v1 - FindUserByID")
		return nil, apperror.ParseGrpc(err)
	}
	return &v1.UserResponse{
		User: toProtoUser(user),
	}, nil
}

func (r *V1) FindUserByEmail(ctx context.Context, req *v1.FindUserByEmailRequest) (*v1.UserResponse, error) {
	user, err := r.authUseCase.FindUserByEmail(ctx, req.Email)
	if err != nil {
		r.l.Error(err, "grpc - v1 - FindUserByEmail")
		return nil, apperror.ParseGrpc(err)
	}
	return &v1.UserResponse{
		User: toProtoUser(user),
	}, nil
}

func (r *V1) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {
	if err := r.authUseCase.Refresh(ctx, req.UserId, req.SessionId, req.NewSessionId, req.ExpiresAt.AsTime()); err != nil {
		r.l.Error(err, "grpc - v1 - RefreshToken")
		return nil, apperror.ParseGrpc(err)
	}
	return &v1.RefreshTokenResponse{
		UserId:       req.UserId,
		NewSessionId: req.NewSessionId,
	}, nil
}
