package v1

import (
	"net/http"

	"github.com/KimNattanan/go-chat-backend/internal/message/entity"
	"github.com/KimNattanan/go-chat-backend/internal/message/handler/rest/v1/request"
	"github.com/KimNattanan/go-chat-backend/pkg/responses"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func (r *V1) findMessageByID(c *echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	message, err := r.messageUseCase.FindByID(ctx, id)
	if err != nil {
		r.l.Error(err, "rest - v1 - findMessageByID")
		return responses.ErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, toMessageResponse(message))
}

func (r *V1) findMessagesByRoomID(c *echo.Context) error {
	ctx := c.Request().Context()
	roomID := c.Param("roomID")

	messages, err := r.messageUseCase.FindByRoomID(ctx, roomID)
	if err != nil {
		r.l.Error(err, "rest - v1 - findMessageByRoomID")
		return responses.ErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, toMessageResponseList(messages))
}

func (r *V1) findMessagesByUserID(c *echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Param("userID")

	messages, err := r.messageUseCase.FindByUserID(ctx, userID)
	if err != nil {
		r.l.Error(err, "rest - v1 - findMessageByUserID")
		return responses.ErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, toMessageResponseList(messages))
}

func (r *V1) findMessagesByRoomIDAndUserID(c *echo.Context) error {
	ctx := c.Request().Context()
	roomID := c.Param("roomID")
	userID := c.Param("userID")

	messages, err := r.messageUseCase.FindByRoomIDAndUserID(ctx, roomID, userID)
	if err != nil {
		r.l.Error(err, "rest - v1 - findMessageByRoomIDAndUserID")
		return responses.ErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, toMessageResponseList(messages))
}

func (r *V1) createMessage(c *echo.Context) error {
	ctx := c.Request().Context()

	roomIDStr := c.Param("roomID")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		r.l.Error(err, "rest - v1 - createMessage")
		return responses.ErrorResponse(c, err)
	}

	userIDStr := c.Get("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		r.l.Error(err, "rest - v1 - createMessage")
		return responses.ErrorResponse(c, err)
	}

	var req request.CreateMessageRequest
	if err := c.Bind(&req); err != nil {
		r.l.Error(err, "rest - v1 - createMessage")
		return responses.ErrorResponse(c, err)
	}
	if err := r.v.Struct(&req); err != nil {
		r.l.Error(err, "rest - v1 - createMessage")
		return responses.ErrorResponse(c, err)
	}

	message := &entity.Message{
		RoomID:  roomID,
		UserID:  userID,
		Content: req.Content,
	}
	if err := r.messageUseCase.Create(ctx, message); err != nil {
		r.l.Error(err, "rest - v1 - createMessage")
		return responses.ErrorResponse(c, err)
	}

	return c.JSON(http.StatusCreated, toMessageResponse(message))
}

func (r *V1) deleteMessage(c *echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	if err := r.messageUseCase.Delete(ctx, id); err != nil {
		r.l.Error(err, "rest - v1 - deleteMessage")
		return responses.ErrorResponse(c, err)
	}

	return responses.MessageResponse(c, http.StatusOK, "message deleted")
}

func (r *V1) deleteMessagesByRoomID(c *echo.Context) error {
	ctx := c.Request().Context()
	roomID := c.Param("roomID")

	if err := r.messageUseCase.DeleteByRoomID(ctx, roomID); err != nil {
		r.l.Error(err, "rest - v1 - deleteMessageByRoomID")
		return responses.ErrorResponse(c, err)
	}

	return responses.MessageResponse(c, http.StatusOK, "messages deleted")
}

func (r *V1) anonymizeUserMessages(c *echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Param("userID")

	if err := r.messageUseCase.AnonymizeUserMessages(ctx, userID); err != nil {
		r.l.Error(err, "rest - v1 - anonymizeUserMessages")
		return responses.ErrorResponse(c, err)
	}

	return responses.MessageResponse(c, http.StatusOK, "messages anonymized")
}
