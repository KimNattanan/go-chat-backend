package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/chat/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
	"github.com/go-playground/validator/v10"
)

// NewChatRoutes -.
func NewChatRoutes(roomUseCase usecase.RoomUseCase, membershipUseCase usecase.MembershipUseCase, l logger.Interface) map[string]rabbitmq.Handler {
	r := &V1{
		roomUseCase: roomUseCase,
		membershipUseCase: membershipUseCase,
		l:              l,
		v:              validator.New(validator.WithRequiredStructEnabled()),
	}
	routes := make(map[string]rabbitmq.Handler)
	routes["user.deleted"] = r.userDeleted
	return routes
}
