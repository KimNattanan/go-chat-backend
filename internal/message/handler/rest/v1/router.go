package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/message/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

// NewMessageRoutes -.
func NewMessageRoutes(apiPrivateGroup *echo.Group, messageUseCase usecase.MessageUseCase, l logger.Interface) {
	r := &V1{
		messageUseCase: messageUseCase,
		l:              l,
		v:              validator.New(validator.WithRequiredStructEnabled()),
	}

	// Private Routes

	messagePrivateGroup := apiPrivateGroup.Group("/messages")
	{
		messagePrivateGroup.GET("/:id", r.findMessageByID)
		messagePrivateGroup.DELETE("/:id", r.deleteMessage)
		messagePrivateGroup.GET("/room/:roomID", r.findMessageByRoomID)
		messagePrivateGroup.DELETE("/room/:roomID", r.deleteMessageByRoomID)
		messagePrivateGroup.POST("/room/:roomID/me", r.createMessage)
		messagePrivateGroup.GET("/room/:roomID/user/:userID", r.findMessageByRoomIDAndUserID)
		messagePrivateGroup.GET("/user/:userID", r.findMessageByUserID)
		messagePrivateGroup.PATCH("/user/:userID/anonymize", r.anonymizeUserMessages)
	}
}
