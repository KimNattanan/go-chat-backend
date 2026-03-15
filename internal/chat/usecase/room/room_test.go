package room

import (
	"context"
	"errors"
	"testing"

	"github.com/KimNattanan/go-chat-backend/internal/chat/entity"
	"github.com/KimNattanan/go-chat-backend/internal/chat/repo"
	"github.com/google/uuid"
)

type mockRoomRepo struct {
	createFunc    func(ctx context.Context, room *entity.Room) error
	findByIDFunc  func(ctx context.Context, id string) (*entity.Room, error)
	findByUserID  func(ctx context.Context, userID string) ([]*entity.Room, error)
	patchFunc     func(ctx context.Context, id string, room *entity.Room) error
	deleteFunc    func(ctx context.Context, id string) error
}

func (m *mockRoomRepo) Create(ctx context.Context, room *entity.Room) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, room)
	}
	return nil
}

func (m *mockRoomRepo) FindByID(ctx context.Context, id string) (*entity.Room, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRoomRepo) FindByUserID(ctx context.Context, userID string) ([]*entity.Room, error) {
	if m.findByUserID != nil {
		return m.findByUserID(ctx, userID)
	}
	return nil, nil
}

func (m *mockRoomRepo) Patch(ctx context.Context, id string, room *entity.Room) error {
	if m.patchFunc != nil {
		return m.patchFunc(ctx, id, room)
	}
	return nil
}

func (m *mockRoomRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

var _ repo.RoomRepo = (*mockRoomRepo)(nil)

type mockEventPublisher struct {
	publishFunc func(msgType string, data any) error
}

func (m *mockEventPublisher) Publish(msgType string, data any) error {
	if m.publishFunc != nil {
		return m.publishFunc(msgType, data)
	}
	return nil
}

func TestRoomUseCase_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var created *entity.Room
		repo := &mockRoomRepo{
			createFunc: func(ctx context.Context, room *entity.Room) error {
				created = room
				return nil
			},
		}
		uc := New(repo, &mockEventPublisher{})

		room := &entity.Room{Title: "test room"}
		err := uc.Create(context.Background(), room)
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		if created == nil || created.Title != "test room" {
			t.Errorf("expected room to be created with Title 'test room', got %+v", created)
		}
	})

	t.Run("repo_error", func(t *testing.T) {
		wantErr := errors.New("db error")
		repo := &mockRoomRepo{
			createFunc: func(context.Context, *entity.Room) error { return wantErr },
		}
		uc := New(repo, &mockEventPublisher{})

		err := uc.Create(context.Background(), &entity.Room{Title: "x"})
		if err != wantErr {
			t.Errorf("Create err = %v, want %v", err, wantErr)
		}
	})
}

func TestRoomUseCase_FindByID(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		id := uuid.New().String()
		want := &entity.Room{Title: "found"}
		repo := &mockRoomRepo{
			findByIDFunc: func(ctx context.Context, rid string) (*entity.Room, error) {
				if rid != id {
					return nil, errors.New("wrong id")
				}
				return want, nil
			},
		}
		uc := New(repo, nil)

		got, err := uc.FindByID(context.Background(), id)
		if err != nil {
			t.Fatalf("FindByID: %v", err)
		}
		if got != want {
			t.Errorf("FindByID = %+v, want %+v", got, want)
		}
	})
}

func TestRoomUseCase_Delete(t *testing.T) {
	t.Run("publish_then_delete", func(t *testing.T) {
		var publishedType string
		var publishedID string
		var deleteCalled bool
		repo := &mockRoomRepo{
			deleteFunc: func(ctx context.Context, id string) error {
				deleteCalled = true
				return nil
			},
		}
		pub := &mockEventPublisher{
			publishFunc: func(msgType string, data any) error {
				publishedType = msgType
				if m, ok := data.(map[string]string); ok {
					publishedID = m["id"]
				}
				return nil
			},
		}
		uc := New(repo, pub)

		id := "room-123"
		err := uc.Delete(context.Background(), id)
		if err != nil {
			t.Fatalf("Delete: %v", err)
		}
		if publishedType != "room.deleted" {
			t.Errorf("Publish msgType = %q, want room.deleted", publishedType)
		}
		if publishedID != id {
			t.Errorf("Publish id = %q, want %q", publishedID, id)
		}
		if !deleteCalled {
			t.Error("Delete was not called on repo")
		}
	})

	t.Run("publish_fails", func(t *testing.T) {
		wantErr := errors.New("publish failed")
		repo := &mockRoomRepo{}
		pub := &mockEventPublisher{publishFunc: func(string, any) error { return wantErr }}
		uc := New(repo, pub)

		err := uc.Delete(context.Background(), "any")
		if err != wantErr {
			t.Errorf("Delete err = %v, want %v", err, wantErr)
		}
	})
}
