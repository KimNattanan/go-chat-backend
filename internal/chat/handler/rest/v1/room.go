package v1

import (
	"net/http"

	"github.com/KimNattanan/go-chat-backend/internal/chat/entity"
	"github.com/KimNattanan/go-chat-backend/internal/chat/handler/rest/v1/request"
	"github.com/KimNattanan/go-chat-backend/pkg/responses"
	"github.com/labstack/echo/v5"
)

func (r *V1) findRoomByID(c *echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	room, err := r.roomUseCase.FindByID(ctx, id)
	if err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - findRoomByID")
	}

	return c.JSON(http.StatusOK, toRoomResponse(room))
}

func (r *V1) findRoomsByUserID(c *echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Get("userID").(string)

	rooms, err := r.roomUseCase.FindByUserID(ctx, userID)
	if err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - findRoomsByUserID")
	}

	return c.JSON(http.StatusOK, toRoomResponseList(rooms))
}

func (r *V1) createRoom(c *echo.Context) error {
	ctx := c.Request().Context()

	var req request.CreateRoomRequest
	if err := c.Bind(&req); err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - createRoom")
	}
	if err := r.v.Struct(&req); err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - createRoom")
	}

	room := &entity.Room{
		Title: req.Title,
	}
	if err := r.roomUseCase.Create(ctx, room); err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - createRoom")
	}

	return c.JSON(http.StatusCreated, toRoomResponse(room))
}

func (r *V1) patchRoom(c *echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	var req request.PatchRoomRequest
	if err := c.Bind(&req); err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - patchRoom")
	}
	if err := r.v.Struct(&req); err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - patchRoom")
	}

	room := &entity.Room{
		Title: req.Title,
	}
	updatedRoom, err := r.roomUseCase.Patch(ctx, id, room)
	if err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - patchRoom")
	}

	return c.JSON(http.StatusOK, toRoomResponse(updatedRoom))
}

func (r *V1) deleteRoom(c *echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	if err := r.roomUseCase.Delete(ctx, id); err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "rest - v1 - deleteRoom")
	}

	return responses.MessageResponse(c, http.StatusOK, "room deleted")
}
