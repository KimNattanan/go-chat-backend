package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KimNattanan/go-chat-backend/internal/auth/entity"
	"github.com/KimNattanan/go-chat-backend/internal/auth/repo"
	"github.com/KimNattanan/go-chat-backend/internal/auth/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/apperror"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
	"github.com/KimNattanan/go-chat-backend/pkg/token"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Precomputed bcrypt hash (cost 12) used when the email is unknown so login
// timing stays closer to the wrong-password path.
var loginDummyPasswordHash = []byte("$2a$12$5OQizdq1EYcgnR7OWfBicO7sXvzEOwTDGTgB71o7AAhoyy2jhtVEK")

type UseCase struct {
	userRepo    repo.UserRepo
	sessionRepo repo.SessionRepo
	mqPublisher rabbitmq.Publisher
	jwtMaker    *token.JWTMaker
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

func New(userRepo repo.UserRepo, sessionRepo repo.SessionRepo, mqPublisher rabbitmq.Publisher, jwtMaker *token.JWTMaker, accessTTL, refreshTTL int) *UseCase {
	return &UseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		mqPublisher: mqPublisher,
		jwtMaker:    jwtMaker,
		accessTTL:   time.Duration(accessTTL) * time.Second,
		refreshTTL:  time.Duration(refreshTTL) * time.Second,
	}
}

func (u *UseCase) FindUserByID(ctx context.Context, id string) (*entity.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

func (u *UseCase) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.userRepo.FindByEmail(ctx, email)
}

func (u *UseCase) DeleteUser(ctx context.Context, id string) error {
	if err := u.sessionRepo.RevokeAllByUserID(ctx, id); err != nil {
		return fmt.Errorf("AuthUseCase - DeleteUser - u.sessionRepo.RevokeAllByUserID: %w", err)
	}
	if err := u.mqPublisher.Publish("user.deleted", map[string]string{
		"user_id": id,
	}); err != nil {
		return fmt.Errorf("AuthUseCase - DeleteUser - u.mqPublisher.Publish: %w", err)
	}
	return u.userRepo.Delete(ctx, id)
}

func (u *UseCase) issueAuthResult(ctx context.Context, user *entity.User, sessionCreatedAt time.Time) (*usecase.AuthResult, error) {
	accessToken, accessClaims, err := u.jwtMaker.CreateToken(user.ID.String(), token.TokenTypeAccess, u.accessTTL)
	if err != nil {
		return nil, fmt.Errorf("issueAuthResult - CreateToken access: %w", err)
	}
	refreshToken, refreshClaims, err := u.jwtMaker.CreateToken(user.ID.String(), token.TokenTypeRefresh, u.refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("issueAuthResult - CreateToken refresh: %w", err)
	}

	sessionID, err := uuid.Parse(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		return nil, fmt.Errorf("issueAuthResult - uuid.Parse: %w", err)
	}
	session := &entity.Session{
		ID:        sessionID,
		UserID:    user.ID,
		IsRevoked: false,
		CreatedAt: sessionCreatedAt,
		ExpiresAt: refreshClaims.ExpiresAt.Time,
	}
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}
	if err := u.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("issueAuthResult - sessionRepo.Create: %w", err)
	}

	return &usecase.AuthResult{
		User:          user,
		AccessToken:   accessToken,
		AccessClaims:  accessClaims,
		RefreshToken:  refreshToken,
		RefreshClaims: refreshClaims,
	}, nil
}

func (u *UseCase) Login(ctx context.Context, email, password string) (*usecase.AuthResult, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = bcrypt.CompareHashAndPassword(loginDummyPasswordHash, []byte(password))
			return nil, apperror.Unauthorized("invalid email or password", err)
		}
		return nil, fmt.Errorf("AuthUseCase - Login - u.userRepo.FindByEmail: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, apperror.Unauthorized("invalid email or password", err)
	}

	result, err := u.issueAuthResult(ctx, user, time.Time{})
	if err != nil {
		return nil, fmt.Errorf("AuthUseCase - Login - %w", err)
	}
	return result, nil
}

func (u *UseCase) Register(ctx context.Context, email, password, name string) (*usecase.AuthResult, error) {
	_, err := u.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, apperror.BadRequest("email already exists", nil)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("AuthUseCase - Register - u.userRepo.FindByEmail: %w", err)
	}
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, fmt.Errorf("AuthUseCase - Register - bcrypt.GenerateFromPassword: %w", err)
	}
	user := &entity.User{
		Email:    email,
		Password: string(hashedPasswordBytes),
	}
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("AuthUseCase - Register - u.userRepo.Create: %w", err)
	}

	if err := u.mqPublisher.Publish("user.created", map[string]string{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"name":    name,
	}); err != nil {
		if delErr := u.userRepo.Delete(ctx, user.ID.String()); delErr != nil {
			return nil, fmt.Errorf("AuthUseCase - Register - u.mqPublisher.Publish: %w (rollback delete: %v)", err, delErr)
		}
		return nil, fmt.Errorf("AuthUseCase - Register - u.mqPublisher.Publish: %w", err)
	}

	result, err := u.issueAuthResult(ctx, user, time.Time{})
	if err != nil {
		return nil, fmt.Errorf("AuthUseCase - Register - %w", err)
	}
	return result, nil
}

func (u *UseCase) Logout(ctx context.Context, refreshToken string) error {
	refreshClaims, err := u.jwtMaker.VerifyToken(refreshToken, token.TokenTypeRefresh)
	if err != nil {
		return fmt.Errorf("AuthUseCase - Logout - u.jwtMaker.VerifyToken: %w", err)
	}
	if err := u.sessionRepo.Delete(ctx, refreshClaims.RegisteredClaims.ID); err != nil {
		return fmt.Errorf("AuthUseCase - Logout - u.sessionRepo.Delete: %w", err)
	}
	return nil
}

func (u *UseCase) RefreshTokenBySessionID(ctx context.Context, userID, oldSessionID string) (*usecase.AuthResult, error) {
	session, err := u.sessionRepo.FindByID(ctx, oldSessionID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, apperror.Unauthorized("invalid refresh token", err)
		}
		return nil, fmt.Errorf("AuthUseCase - RefreshTokenBySessionID - u.sessionRepo.FindByID: %w", err)
	}
	if session.UserID.String() != userID {
		return nil, apperror.Unauthorized("invalid refresh token", fmt.Errorf("AuthUseCase - RefreshTokenBySessionID: session user ID does not match"))
	}
	if session.IsRevoked {
		return nil, apperror.Unauthorized("invalid refresh token", nil)
	}

	if err := u.sessionRepo.Revoke(ctx, oldSessionID); err != nil {
		return nil, fmt.Errorf("AuthUseCase - RefreshTokenBySessionID - u.sessionRepo.Revoke: %w", err)
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("AuthUseCase - RefreshTokenBySessionID - u.userRepo.FindByID: %w", err)
	}

	result, err := u.issueAuthResult(ctx, user, session.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("AuthUseCase - RefreshTokenBySessionID - %w", err)
	}
	return result, nil
}
