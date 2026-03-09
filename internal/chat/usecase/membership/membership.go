package membership

import (
	"context"

	authPb "github.com/KimNattanan/go-chat-backend/internal/auth/proto/v1"
	"github.com/KimNattanan/go-chat-backend/internal/chat/entity"
	"github.com/KimNattanan/go-chat-backend/internal/chat/repo"
)

type UseCase struct {
	membershipRepo repo.MembershipRepo
	authGrpcClient authPb.AuthServiceClient
}

func New(membershipRepo repo.MembershipRepo, authGrpcClient authPb.AuthServiceClient) *UseCase {
	return &UseCase{
		membershipRepo: membershipRepo,
		authGrpcClient: authGrpcClient,
	}
}

func (u *UseCase) Create(ctx context.Context, membership *entity.Membership) error {
	_, err := u.authGrpcClient.FindUserByID(ctx, &authPb.FindUserByIDRequest{
		Id: membership.UserID.String(),
	})
	if err != nil {
		return err
	}
	return u.membershipRepo.Create(ctx, membership)
}

func (u *UseCase) FindByRoomID(ctx context.Context, roomID string) ([]*entity.Membership, error) {
	return u.membershipRepo.FindByRoomID(ctx, roomID)
}

func (u *UseCase) FindByUserID(ctx context.Context, userID string) ([]*entity.Membership, error) {
	return u.membershipRepo.FindByUserID(ctx, userID)
}

func (u *UseCase) FindByRoomIDAndUserID(ctx context.Context, roomID, userID string) (*entity.Membership, error) {
	return u.membershipRepo.FindByRoomIDAndUserID(ctx, roomID, userID)
}

func (u *UseCase) Delete(ctx context.Context, roomID, userID string) error {
	return u.membershipRepo.Delete(ctx, roomID, userID)
}

func (u *UseCase) DeleteByUserID(ctx context.Context, userID string) error {
	return u.membershipRepo.DeleteByUserID(ctx, userID)
}
