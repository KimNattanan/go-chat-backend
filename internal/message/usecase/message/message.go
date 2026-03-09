package message

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/message/entity"
	"github.com/KimNattanan/go-chat-backend/internal/message/repo"
)

type UseCase struct {
	messageRepo repo.MessageRepo
}

func New(messageRepo repo.MessageRepo) *UseCase {
	return &UseCase{
		messageRepo: messageRepo,
	}
}

func (u *UseCase) Create(ctx context.Context, message *entity.Message) error {
	return u.messageRepo.Create(ctx, message)
}

func (u *UseCase) FindByID(ctx context.Context, id string) (*entity.Message, error) {
	return u.messageRepo.FindByID(ctx, id)
}

func (u *UseCase) FindByRoomID(ctx context.Context, roomID string) ([]*entity.Message, error) {
	return u.messageRepo.FindByRoomID(ctx, roomID)
}

func (u *UseCase) FindByUserID(ctx context.Context, userID string) ([]*entity.Message, error) {
	return u.messageRepo.FindByUserID(ctx, userID)
}

func (u *UseCase) AnonymizeUserMessages(ctx context.Context, userID string) error {
	return u.messageRepo.AnonymizeUserMessages(ctx, userID)
}

func (u *UseCase) Delete(ctx context.Context, id string) error {
	return u.messageRepo.Delete(ctx, id)
}

func (u *UseCase) DeleteByRoomID(ctx context.Context, roomID string) error {
	return u.messageRepo.DeleteByRoomID(ctx, roomID)
}
