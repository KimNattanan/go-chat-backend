package usecase

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/auth/entity"
	"github.com/KimNattanan/go-chat-backend/pkg/token"
)

type AuthResult struct {
	User          *entity.User
	AccessToken   string
	AccessClaims  *token.UserClaims
	RefreshToken  string
	RefreshClaims *token.UserClaims
}

type (
	AuthUseCase interface {
		FindUserByID(ctx context.Context, id string) (*entity.User, error)
		FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
		DeleteUser(ctx context.Context, id string) error

		Login(ctx context.Context, email, password string) (*AuthResult, error)
		Register(ctx context.Context, email, password, name string) (*AuthResult, error)
		Logout(ctx context.Context, accessToken, refreshToken string) error
		RefreshTokenBySessionID(ctx context.Context, userID, oldSessionID string) (*AuthResult, error)
	}
)
