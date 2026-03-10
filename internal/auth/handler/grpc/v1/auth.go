package v1

import (
	"context"

	v1 "github.com/KimNattanan/go-chat-backend/internal/auth/proto/v1"
	"github.com/KimNattanan/go-chat-backend/pkg/apperror"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *V1) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	user, accessToken, accessClaims, refreshToken, refreshClaims, err := r.authUseCase.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.LoginResponse{
		UserId: user.ID.String(),
		Tokens: &v1.Tokens{
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  timestamppb.New(accessClaims.ExpiresAt.Time),
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: timestamppb.New(refreshClaims.ExpiresAt.Time),
		},
	}, nil
}

func (r *V1) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	user, accessToken, accessClaims, refreshToken, refreshClaims, err := r.authUseCase.Register(ctx, req.Email, req.Password, req.Profile.Name)
	if err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.RegisterResponse{
		UserId: user.ID.String(),
		Tokens: &v1.Tokens{
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  timestamppb.New(accessClaims.ExpiresAt.Time),
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: timestamppb.New(refreshClaims.ExpiresAt.Time),
		},
	}, nil
}

func (r *V1) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	if err := r.authUseCase.Logout(ctx, req.RefreshToken); err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.LogoutResponse{
		Message: "logged out successfully",
	}, nil
}

func (r *V1) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {
	user, accessToken, accessClaims, refreshToken, refreshClaims, err := r.authUseCase.RefreshTokenBySessionID(ctx, req.UserId, req.SessionId)
	if err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.RefreshTokenResponse{
		UserId: user.ID.String(),
		Tokens: &v1.Tokens{
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  timestamppb.New(accessClaims.RegisteredClaims.ExpiresAt.Time),
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: timestamppb.New(refreshClaims.RegisteredClaims.ExpiresAt.Time),
		},
	}, nil
}

func (r *V1) FindUserByID(ctx context.Context, req *v1.FindUserByIDRequest) (*v1.UserResponse, error) {
	user, err := r.authUseCase.FindUserByID(ctx, req.Id)
	if err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.UserResponse{
		User: toProtoUser(user),
	}, nil
}

func (r *V1) FindUserByEmail(ctx context.Context, req *v1.FindUserByEmailRequest) (*v1.UserResponse, error) {
	user, err := r.authUseCase.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.UserResponse{
		User: toProtoUser(user),
	}, nil
}

func (r *V1) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error) {
	if err := r.authUseCase.DeleteUser(ctx, req.Id); err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.DeleteUserResponse{
		Message: "user deleted",
	}, nil
}
