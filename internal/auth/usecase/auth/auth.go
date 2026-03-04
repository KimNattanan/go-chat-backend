package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KimNattanan/go-chat-backend/internal/auth/entity"
	"github.com/KimNattanan/go-chat-backend/internal/auth/repo"
	profilePb "github.com/KimNattanan/go-chat-backend/internal/profile/proto/v1"
	"github.com/KimNattanan/go-chat-backend/pkg/token"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UseCase struct {
	userRepo          repo.UserRepo
	sessionRepo       repo.SessionRepo
	profileGrpcClient profilePb.ProfileServiceClient
	jwtMaker          *token.JWTMaker
	accessTTL         time.Duration
	refreshTTL        time.Duration
}

func New(userRepo repo.UserRepo, sessionRepo repo.SessionRepo, profileGrpcClient profilePb.ProfileServiceClient, jwtMaker *token.JWTMaker, accessTTL, refreshTTL int) *UseCase {
	return &UseCase{
		userRepo:          userRepo,
		sessionRepo:       sessionRepo,
		profileGrpcClient: profileGrpcClient,
		jwtMaker:          jwtMaker,
		accessTTL:         time.Duration(accessTTL),
		refreshTTL:        time.Duration(refreshTTL),
	}
}

func (u *UseCase) FindUserByID(ctx context.Context, id string) (*entity.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

func (u *UseCase) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.userRepo.FindByEmail(ctx, email)
}

func (u *UseCase) DeleteUser(ctx context.Context, id string) error {
	if _, err := u.profileGrpcClient.DeleteProfile(ctx, &profilePb.DeleteProfileRequest{
		UserId: id,
	}); err != nil {
		return fmt.Errorf("AuthUseCase - DeleteUser - u.profileGrpcClient.DeleteProfile: %w", err)
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
	if _, err := u.profileGrpcClient.CreateProfile(ctx, &profilePb.CreateProfileRequest{
		UserId: user.ID.String(),
		Email:  user.Email,
		Name:   name,
	}); err != nil {
		u.userRepo.Delete(ctx, user.ID.String())
		return nil, "", nil, "", nil, fmt.Errorf("AuthUseCase - Register - u.profileGrpcClient.CreateProfile: %w", err)
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

func (u *UseCase) Refresh(ctx context.Context, userID, sessionID, newIDStr string, expiresAt time.Time) error {
	newID, err := uuid.Parse(newIDStr)
	if err != nil {
		return fmt.Errorf("AuthUseCase - Refresh - uuid.Parse: %w", err)
	}
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("AuthUseCase - Refresh - u.userRepo.FindByID: %w", err)
	}
	session, err := u.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("AuthUseCase - Refresh - u.sessionRepo.FindByID: %w", err)
	}
	if session.IsRevoked {
		return fmt.Errorf("AuthUseCase - Refresh: session is revoked")
	}
	if err := u.sessionRepo.Revoke(ctx, sessionID); err != nil {
		return fmt.Errorf("AuthUseCase - Refresh - u.sessionRepo.Revoke: %w", err)
	}
	newSession := &entity.Session{
		ID:        newID,
		UserID:    user.ID,
		IsRevoked: false,
		CreatedAt: session.CreatedAt,
		ExpiresAt: expiresAt,
	}
	if err := u.sessionRepo.Create(ctx, newSession); err != nil {
		return fmt.Errorf("AuthUseCase - Refresh - u.sessionRepo.Create: %w", err)
	}
	return nil
}
