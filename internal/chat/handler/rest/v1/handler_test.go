package v1

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/KimNattanan/go-chat-backend/internal/chat/entity"
	"github.com/KimNattanan/go-chat-backend/internal/chat/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// noopLogger สำหรับเทส ไม่ log อะไร
type noopLogger struct{}

func (noopLogger) Debug(message any, args ...any) {}
func (noopLogger) Info(message string, args ...any) {}
func (noopLogger) Warn(message string, args ...any) {}
func (noopLogger) Error(message any, args ...any) {}
func (noopLogger) Fatal(message any, args ...any) {}

var _ logger.Interface = (*noopLogger)(nil)

// fakeRoomUseCase ใช้ในเทส
type fakeRoomUseCase struct {
	createFunc    func(ctx context.Context, room *entity.Room) error
	findByIDFunc  func(ctx context.Context, id string) (*entity.Room, error)
	findByUserID  func(ctx context.Context, userID string) ([]*entity.Room, error)
	patchFunc     func(ctx context.Context, id string, room *entity.Room) (*entity.Room, error)
	deleteFunc    func(ctx context.Context, id string) error
}

func (f *fakeRoomUseCase) Create(ctx context.Context, room *entity.Room) error {
	if f.createFunc != nil {
		return f.createFunc(ctx, room)
	}
	return nil
}

func (f *fakeRoomUseCase) FindByID(ctx context.Context, id string) (*entity.Room, error) {
	if f.findByIDFunc != nil {
		return f.findByIDFunc(ctx, id)
	}
	return nil, nil
}

func (f *fakeRoomUseCase) FindByUserID(ctx context.Context, userID string) ([]*entity.Room, error) {
	if f.findByUserID != nil {
		return f.findByUserID(ctx, userID)
	}
	return nil, nil
}

func (f *fakeRoomUseCase) Patch(ctx context.Context, id string, room *entity.Room) (*entity.Room, error) {
	if f.patchFunc != nil {
		return f.patchFunc(ctx, id, room)
	}
	return nil, nil
}

func (f *fakeRoomUseCase) Delete(ctx context.Context, id string) error {
	if f.deleteFunc != nil {
		return f.deleteFunc(ctx, id)
	}
	return nil
}

var _ usecase.RoomUseCase = (*fakeRoomUseCase)(nil)

// fakeMembershipUseCase ใช้ในเทส
type fakeMembershipUseCase struct {
	createFunc              func(ctx context.Context, m *entity.Membership) error
	findByRoomIDFunc        func(ctx context.Context, roomID string) ([]*entity.Membership, error)
	findByUserIDFunc        func(ctx context.Context, userID string) ([]*entity.Membership, error)
	findByRoomIDAndUserID   func(ctx context.Context, roomID, userID string) (*entity.Membership, error)
	deleteFunc              func(ctx context.Context, roomID, userID string) error
	deleteByUserIDFunc      func(ctx context.Context, userID string) error
}

func (f *fakeMembershipUseCase) Create(ctx context.Context, m *entity.Membership) error {
	if f.createFunc != nil {
		return f.createFunc(ctx, m)
	}
	return nil
}

func (f *fakeMembershipUseCase) FindByRoomID(ctx context.Context, roomID string) ([]*entity.Membership, error) {
	if f.findByRoomIDFunc != nil {
		return f.findByRoomIDFunc(ctx, roomID)
	}
	return nil, nil
}

func (f *fakeMembershipUseCase) FindByUserID(ctx context.Context, userID string) ([]*entity.Membership, error) {
	if f.findByUserIDFunc != nil {
		return f.findByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (f *fakeMembershipUseCase) FindByRoomIDAndUserID(ctx context.Context, roomID, userID string) (*entity.Membership, error) {
	if f.findByRoomIDAndUserID != nil {
		return f.findByRoomIDAndUserID(ctx, roomID, userID)
	}
	return nil, nil
}

func (f *fakeMembershipUseCase) Delete(ctx context.Context, roomID, userID string) error {
	if f.deleteFunc != nil {
		return f.deleteFunc(ctx, roomID, userID)
	}
	return nil
}

func (f *fakeMembershipUseCase) DeleteByUserID(ctx context.Context, userID string) error {
	if f.deleteByUserIDFunc != nil {
		return f.deleteByUserIDFunc(ctx, userID)
	}
	return nil
}

var _ usecase.MembershipUseCase = (*fakeMembershipUseCase)(nil)

func newTestHandler(roomUC usecase.RoomUseCase, membershipUC usecase.MembershipUseCase) *V1 {
	return &V1{
		roomUseCase:       roomUC,
		membershipUseCase: membershipUC,
		l:                 noopLogger{},
		v:                 mustValidator(),
	}
}

func mustValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	return v
}

func TestV1_findRoomByID(t *testing.T) {
	e := echo.New()
	h := newTestHandler(
		&fakeRoomUseCase{
			findByIDFunc: func(ctx context.Context, id string) (*entity.Room, error) {
				if id == "not-found" {
					return nil, gorm.ErrRecordNotFound
				}
				return &entity.Room{
					ID:        uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					Title:     "test room",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
		},
		&fakeMembershipUseCase{},
	)

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:id")
		c.SetPathValues(echo.PathValues{{Name: "id", Value: "550e8400-e29b-41d4-a716-446655440000"}})

		err := h.findRoomByID(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "test room")
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/rooms/:id")
		c.SetPathValues(echo.PathValues{{Name: "id", Value: "not-found"}})

		err := h.findRoomByID(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestV1_findRoomsByUserID(t *testing.T) {
	e := echo.New()
	userID := "user-123"
	h := newTestHandler(
		&fakeRoomUseCase{
			findByUserID: func(ctx context.Context, uid string) ([]*entity.Room, error) {
				assert.Equal(t, userID, uid)
				return []*entity.Room{
					{ID: uuid.New(), Title: "room1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
				}, nil
			},
		},
		&fakeMembershipUseCase{},
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	err := h.findRoomsByUserID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "room1")
}

func TestV1_createRoom(t *testing.T) {
	e := echo.New()
	var created *entity.Room
	h := newTestHandler(
		&fakeRoomUseCase{
			createFunc: func(ctx context.Context, room *entity.Room) error {
				created = room
				room.ID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				room.CreatedAt = time.Now()
				room.UpdatedAt = time.Now()
				return nil
			},
		},
		&fakeMembershipUseCase{},
	)

	body := `{"title":"new room"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.createRoom(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	require.NotNil(t, created)
	assert.Equal(t, "new room", created.Title)
}

func TestV1_patchRoom(t *testing.T) {
	e := echo.New()
	roomID := "550e8400-e29b-41d4-a716-446655440000"
	h := newTestHandler(
		&fakeRoomUseCase{
			patchFunc: func(ctx context.Context, id string, room *entity.Room) (*entity.Room, error) {
				assert.Equal(t, roomID, id)
				assert.Equal(t, "updated title", room.Title)
				return &entity.Room{
					ID:        uuid.MustParse(roomID),
					Title:     "updated title",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
		},
		&fakeMembershipUseCase{},
	)

	body := `{"title":"updated title"}`
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/rooms/:id")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: roomID}})

	err := h.patchRoom(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "updated title")
}

func TestV1_deleteRoom(t *testing.T) {
	e := echo.New()
	var deletedID string
	h := newTestHandler(
		&fakeRoomUseCase{
			deleteFunc: func(ctx context.Context, id string) error {
				deletedID = id
				return nil
			},
		},
		&fakeMembershipUseCase{},
	)

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/rooms/:id")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "room-to-delete"}})

	err := h.deleteRoom(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "room-to-delete", deletedID)
}

func TestV1_findMembershipsByRoomID(t *testing.T) {
	e := echo.New()
	roomID := uuid.New().String()
	h := newTestHandler(
		&fakeRoomUseCase{},
		&fakeMembershipUseCase{
			findByRoomIDFunc: func(ctx context.Context, rid string) ([]*entity.Membership, error) {
				assert.Equal(t, roomID, rid)
				return []*entity.Membership{
					{RoomID: uuid.MustParse(roomID), UserID: uuid.New(), CreatedAt: time.Now()},
				}, nil
			},
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/memberships/room/:roomID")
	c.SetPathValues(echo.PathValues{{Name: "roomID", Value: roomID}})

	err := h.findMembershipsByRoomID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestV1_findMembershipByRoomIDAndUserID(t *testing.T) {
	e := echo.New()
	roomID := uuid.New().String()
	userID := uuid.New().String()
	h := newTestHandler(
		&fakeRoomUseCase{},
		&fakeMembershipUseCase{
			findByRoomIDAndUserID: func(ctx context.Context, rid, uid string) (*entity.Membership, error) {
				assert.Equal(t, roomID, rid)
				assert.Equal(t, userID, uid)
				return &entity.Membership{
					RoomID: uuid.MustParse(roomID), UserID: uuid.MustParse(userID), CreatedAt: time.Now(),
				}, nil
			},
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/memberships/room/:roomID/user/:userID")
	c.SetPathValues(echo.PathValues{
		{Name: "roomID", Value: roomID},
		{Name: "userID", Value: userID},
	})

	err := h.findMembershipByRoomIDAndUserID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestV1_createMembership(t *testing.T) {
	e := echo.New()
	roomID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	userID := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	var created *entity.Membership
	h := newTestHandler(
		&fakeRoomUseCase{},
		&fakeMembershipUseCase{
			createFunc: func(ctx context.Context, m *entity.Membership) error {
				created = m
				return nil
			},
		},
	)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/memberships/room/:roomID/user/:userID")
	c.SetPathValues(echo.PathValues{
		{Name: "roomID", Value: roomID.String()},
		{Name: "userID", Value: userID.String()},
	})

	err := h.createMembership(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	require.NotNil(t, created)
	assert.Equal(t, roomID, created.RoomID)
	assert.Equal(t, userID, created.UserID)
}

func TestV1_deleteMembershipByRoomIDAndUserID(t *testing.T) {
	e := echo.New()
	roomID := "room-1"
	userID := "user-1"
	var deletedRoomID, deletedUserID string
	h := newTestHandler(
		&fakeRoomUseCase{},
		&fakeMembershipUseCase{
			deleteFunc: func(ctx context.Context, rid, uid string) error {
				deletedRoomID, deletedUserID = rid, uid
				return nil
			},
		},
	)

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/memberships/room/:roomID/user/:userID")
	c.SetPathValues(echo.PathValues{
		{Name: "roomID", Value: roomID},
		{Name: "userID", Value: userID},
	})

	err := h.deleteMembershipByRoomIDAndUserID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, roomID, deletedRoomID)
	assert.Equal(t, userID, deletedUserID)
}

func TestV1_findMembershipsByUserID(t *testing.T) {
	e := echo.New()
	userID := "user-123"
	h := newTestHandler(
		&fakeRoomUseCase{},
		&fakeMembershipUseCase{
			findByRoomIDFunc: func(ctx context.Context, roomID string) ([]*entity.Membership, error) {
				assert.Equal(t, userID, roomID)
				return []*entity.Membership{}, nil
			},
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/memberships/user/:userID")
	c.SetPathValues(echo.PathValues{{Name: "userID", Value: userID}})

	err := h.findMembershipsByUserID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestV1_deleteMembershipsByUserID(t *testing.T) {
	e := echo.New()
	userID := "user-to-delete"
	var deletedUserID string
	h := newTestHandler(
		&fakeRoomUseCase{},
		&fakeMembershipUseCase{
			deleteByUserIDFunc: func(ctx context.Context, uid string) error {
				deletedUserID = uid
				return nil
			},
		},
	)

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/memberships/user/:userID")
	c.SetPathValues(echo.PathValues{{Name: "userID", Value: userID}})

	err := h.deleteMembershipsByUserID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, userID, deletedUserID)
}

func TestV1_createRoom_bindError(t *testing.T) {
	e := echo.New()
	h := newTestHandler(&fakeRoomUseCase{}, &fakeMembershipUseCase{})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.createRoom(c)
	require.NoError(t, err)
	assert.True(t, rec.Code >= 400)
}

func TestV1_findRoomByID_useCaseError(t *testing.T) {
	e := echo.New()
	h := newTestHandler(
		&fakeRoomUseCase{
			findByIDFunc: func(ctx context.Context, id string) (*entity.Room, error) {
				return nil, errors.New("db error")
			},
		},
		&fakeMembershipUseCase{},
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/rooms/:id")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "any"}})

	err := h.findRoomByID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
