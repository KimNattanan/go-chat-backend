package room

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/chat/entity"
	"github.com/KimNattanan/go-chat-backend/internal/chat/repo"
)

type eventPublisher interface {
	Publish(msgType string, data any) error
}

type UseCase struct {
	roomRepo       repo.RoomRepo
	eventPublisher eventPublisher
}

func New(roomRepo repo.RoomRepo, eventPublisher eventPublisher) *UseCase {
	return &UseCase{
		roomRepo:       roomRepo,
		eventPublisher: eventPublisher,
	}
}

func (u *UseCase) Create(ctx context.Context, room *entity.Room) error {
	return u.roomRepo.Create(ctx, room)
}

func (u *UseCase) FindByID(ctx context.Context, id string) (*entity.Room, error) {
	return u.roomRepo.FindByID(ctx, id)
}

func (u *UseCase) FindByUserID(ctx context.Context, userID string) ([]*entity.Room, error) {
	return u.roomRepo.FindByUserID(ctx, userID)
}

func (u *UseCase) Patch(ctx context.Context, id string, room *entity.Room) (*entity.Room, error) {
	if err := u.roomRepo.Patch(ctx, id, room); err != nil {
		return &entity.Room{}, err
	}
	return u.roomRepo.FindByID(ctx, id)
}

func (u *UseCase) Delete(ctx context.Context, id string) error {
	if err := u.eventPublisher.Publish("room.deleted", map[string]string{
		"id": id,
	}); err != nil {
		return err
	}
	return u.roomRepo.Delete(ctx, id)
}
