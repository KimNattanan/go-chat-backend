package usecase

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/auth/entity"
	"github.com/KimNattanan/go-chat-backend/pkg/token"
)

type (
	AuthUseCase interface {
		FindUserByID(ctx context.Context, id string) (*entity.User, error)
		FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
		DeleteUser(ctx context.Context, id string) error

		CreateSession(ctx context.Context, session *entity.Session) error
		FindSessionByID(ctx context.Context, id string) (*entity.Session, error)
		FindSessionByUserID(ctx context.Context, userID string) ([]*entity.Session, error)
		RevokeSession(ctx context.Context, id string) error
		DeleteSession(ctx context.Context, id string) error

		Login(ctx context.Context, email, password string) (*entity.User, string, *token.UserClaims, string, *token.UserClaims, error)
		Register(ctx context.Context, email, password, name string) (*entity.User, string, *token.UserClaims, string, *token.UserClaims, error)
		Logout(ctx context.Context, refreshToken string) error
		RefreshToken(ctx context.Context, userID, oldRefreshToken string) (*entity.User, string, *token.UserClaims, string, *token.UserClaims, error)
		RefreshTokenBySessionID(ctx context.Context, userID, oldSessionID string) (*entity.User, string, *token.UserClaims, string, *token.UserClaims, error)
	}
)
