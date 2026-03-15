package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KimNattanan/go-chat-backend/internal/auth/entity"
	"github.com/KimNattanan/go-chat-backend/internal/auth/repo"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
	"github.com/KimNattanan/go-chat-backend/pkg/token"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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
	if err := u.mqPublisher.Publish("user.deleted", map[string]string{
		"user_id": id,
	}); err != nil {
		return fmt.Errorf("AuthUseCase - DeleteUser - u.mqPublisher.Publish: %w", err)
	}
	return u.userRepo.Delete(ctx, id)
}

func (u *UseCase) CreateSession(ctx context.Context, session *entity.Session) error {
	return u.sessionRepo.Create(ctx, session)
}

func (u *UseCase) FindSessionByID(ctx context.Context, id string) (*entity.Session, error) {
	return u.sessionRepo.FindByID(ctx, id)
}

func (u *UseCase) FindSessionByUserID(ctx context.Context, userID string) ([]*entity.Session, error) {
	return u.sessionRepo.FindByUserID(ctx, userID)
}

func (u *UseCase) RevokeSession(ctx context.Context, id string) error {
	return u.sessionRepo.Revoke(ctx, id)
}

func (u *UseCase) DeleteSession(ctx context.Context, id string) error {
	return u.sessionRepo.Delete(ctx, id)
}

func (u *UseCase) Login(ctx context.Context, email, password string) (*entity.User, string, *token.UserClaims, string, *token.UserClaims, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Login - u.userRepo.FindByEmail: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Login - bcrypt.CompareHashAndPassword: %w", err)
	}

	accessToken, accessClaims, err := u.jwtMaker.CreateToken(user.ID.String(), time.Second*u.accessTTL)
	if err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Login - u.jwtMaker.CreateToken: %w", err)
	}
	refreshToken, refreshClaims, err := u.jwtMaker.CreateToken(user.ID.String(), time.Second*u.refreshTTL)
	if err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Login - u.jwtMaker.CreateToken: %w", err)
	}

	sessionID, err := uuid.Parse(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Login - uuid.Parse: %w", err)
	}
	session := &entity.Session{
		ID:        sessionID,
		UserID:    user.ID,
		IsRevoked: false,
		ExpiresAt: refreshClaims.ExpiresAt.Time,
	}
	if err := u.sessionRepo.Create(ctx, session); err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Login - u.sessionRepo.Create: %w", err)
	}

	return user, accessToken, accessClaims, refreshToken, refreshClaims, nil
}

func (u *UseCase) Register(ctx context.Context, email, password, name string) (*entity.User, string, *token.UserClaims, string, *token.UserClaims, error) {
	_, err := u.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Register: email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Register - u.userRepo.FindByEmail: %w", err)
	}
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Register - bcrypt.GenerateFromPassword: %w", err)
	}
	user := &entity.User{
		Email:    email,
		Password: string(hashedPasswordBytes),
	}
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Register - u.userRepo.Create: %w", err)
	}

	if err := u.mqPublisher.Publish("user.created", map[string]string{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"name":    name,
	}); err != nil {
		u.userRepo.Delete(ctx, user.ID.String())
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Register - u.amqpClient.Publish: %w", err)
	}

	accessToken, accessClaims, err := u.jwtMaker.CreateToken(user.ID.String(), time.Second*u.accessTTL)
	if err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Register - u.jwtMaker.CreateToken: %w", err)
	}
	refreshToken, refreshClaims, err := u.jwtMaker.CreateToken(user.ID.String(), time.Second*u.refreshTTL)
	if err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Register - u.jwtMaker.CreateToken: %w", err)
	}

	sessionID, err := uuid.Parse(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Register - uuid.Parse: %w", err)
	}
	session := &entity.Session{
		ID:        sessionID,
		UserID:    user.ID,
		IsRevoked: false,
		ExpiresAt: refreshClaims.ExpiresAt.Time,
	}
	if err := u.sessionRepo.Create(ctx, session); err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Register - u.sessionRepo.Create: %w", err)
	}

	return user, accessToken, accessClaims, refreshToken, refreshClaims, nil
}

func (u *UseCase) Logout(ctx context.Context, refreshToken string) error {
	refreshClaims, err := u.jwtMaker.VerifyToken(refreshToken)
	if err != nil {
		return fmt.Errorf("AuthUseCase - Logout - u.jwtMaker.VerifyToken: %w", err)
	}
	if err := u.sessionRepo.Revoke(ctx, refreshClaims.RegisteredClaims.ID); err != nil {
		return fmt.Errorf("AuthUseCase - Logout - u.sessionRepo.Revoke: %w", err)
	}
	return nil
}

func (u *UseCase) RefreshTokenBySessionID(ctx context.Context, userID, oldSessionID string) (*entity.User, string, *token.UserClaims, string, *token.UserClaims, error) {
	session, err := u.sessionRepo.FindByID(ctx, oldSessionID)
	if err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - RefreshTokenBySessionID - u.sessionRepo.FindByID: %w", err)
	}
	if session.IsRevoked {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - RefreshTokenBySessionID: session is revoked")
	}
	if err := u.sessionRepo.Revoke(ctx, oldSessionID); err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - RefreshTokenBySessionID - u.sessionRepo.Revoke: %w", err)
	}
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - RefreshTokenBySessionID - u.userRepo.FindByID: %w", err)
	}

	accessToken, accessClaims, err := u.jwtMaker.CreateToken(user.ID.String(), time.Second*u.accessTTL)
	if err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - RefreshTokenBySessionID - u.jwtMaker.CreateToken: %w", err)
	}
	refreshToken, refreshClaims, err := u.jwtMaker.CreateToken(user.ID.String(), time.Second*u.refreshTTL)
	if err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - RefreshTokenBySessionID - u.jwtMaker.CreateToken: %w", err)
	}

	newSessionID, err := uuid.Parse(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - RefreshTokenBySessionID - uuid.Parse: %w", err)
	}
	newSession := &entity.Session{
		ID:        newSessionID,
		UserID:    user.ID,
		IsRevoked: false,
		CreatedAt: session.CreatedAt,
		ExpiresAt: refreshClaims.ExpiresAt.Time,
	}
	if err := u.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - RefreshTokenBySessionID - u.sessionRepo.Create: %w", err)
	}

	return user, accessToken, accessClaims, refreshToken, refreshClaims, nil
}

func (u *UseCase) RefreshToken(ctx context.Context, userID, oldRefreshToken string) (*entity.User, string, *token.UserClaims, string, *token.UserClaims, error) {
	oldRefreshClaims, err := u.jwtMaker.VerifyToken(oldRefreshToken)
	if err != nil {
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - RefreshToken - u.jwtMaker.VerifyToken: %w", err)
	}
	return u.RefreshTokenBySessionID(ctx, userID, oldRefreshClaims.RegisteredClaims.ID)
}
