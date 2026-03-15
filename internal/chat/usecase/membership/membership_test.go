package membership_test

import (
	"context"
	"errors"
	"testing"

	authPb "github.com/KimNattanan/go-chat-backend/internal/auth/proto/v1"
	authMocks "github.com/KimNattanan/go-chat-backend/internal/auth/proto/v1/mocks"
	"github.com/KimNattanan/go-chat-backend/internal/chat/entity"
	"github.com/KimNattanan/go-chat-backend/internal/chat/repo"
	membershipUsecase "github.com/KimNattanan/go-chat-backend/internal/chat/usecase/membership"
	"github.com/KimNattanan/go-chat-backend/internal/chat/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var _ repo.MembershipRepo = (*mocks.MockMembershipRepo)(nil)

func TestUseCase_Create(t *testing.T) {
	t.Run("success when auth client finds user", func(t *testing.T) {
		userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		authClient := authMocks.NewMockAuthServiceClient(t)
		authClient.EXPECT().
			FindUserByID(context.Background(), &authPb.FindUserByIDRequest{Id: userID.String()}).
			Return(&authPb.UserResponse{}, nil)

		membership := &entity.Membership{
			UserID: userID,
			RoomID: uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
		}
		
		membershipRepo := mocks.NewMockMembershipRepo(t)
		membershipRepo.EXPECT().
			Create(context.Background(), membership).
			Once().
			Return(nil)
		
		uc := membershipUsecase.New(membershipRepo, authClient)
		err := uc.Create(context.Background(), membership)
		assert.NoError(t, err)
	})

	t.Run("error when auth client returns error", func(t *testing.T) {
		userID := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		authClient := authMocks.NewMockAuthServiceClient(t)
		authClient.EXPECT().
			FindUserByID(context.Background(), &authPb.FindUserByIDRequest{Id: userID.String()}).
			Return(nil, errors.New("user not found"))

		membershipRepo := mocks.NewMockMembershipRepo(t)
		uc := membershipUsecase.New(membershipRepo, authClient)
		roomID := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		membership := &entity.Membership{UserID: userID, RoomID: roomID}

		err := uc.Create(context.Background(), membership)
		if err == nil {
			t.Fatal("expected error from Create when auth fails")
		}
	})
}
