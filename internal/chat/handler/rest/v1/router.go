package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/chat/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

// NewChatRoutes -.
func NewChatRoutes(apiPrivateGroup *echo.Group, roomUseCase usecase.RoomUseCase, membershipUseCase usecase.MembershipUseCase, l logger.Interface) {
	r := &V1{
		roomUseCase:       roomUseCase,
		membershipUseCase: membershipUseCase,
		l:                 l,
		v:                 validator.New(validator.WithRequiredStructEnabled()),
	}

	// Private Routes

	roomPrivateGroup := apiPrivateGroup.Group("/rooms")
	{
		roomPrivateGroup.GET("/:id", r.findRoomByID)
		roomPrivateGroup.POST("", r.createRoom)
		roomPrivateGroup.PATCH("/:id", r.patchRoom)
		roomPrivateGroup.DELETE("/:id", r.deleteRoom)
	}
	membershipPrivateGroup := apiPrivateGroup.Group("/memberships")
	{
		membershipPrivateGroup.GET("/room/:roomID", r.findMembershipsByRoomID)
		membershipPrivateGroup.GET("/room/:roomID/user/:userID", r.findMembershipByRoomIDAndUserID)
		membershipPrivateGroup.POST("/room/:roomID/user/:userID", r.createMembership)
		membershipPrivateGroup.DELETE("/room/:roomID/user/:userID", r.deleteMembershipByRoomIDAndUserID)
		membershipPrivateGroup.GET("/user/:userID", r.findMembershipsByUserID)
		membershipPrivateGroup.DELETE("/user/:userID", r.deleteMembershipsByUserID)
	}

	userPrivateGroup := apiPrivateGroup.Group("/users")
	{
		userPrivateGroup.GET("/rooms", r.findRoomsByUserID)
	}
}
