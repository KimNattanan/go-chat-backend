package v1

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/chat/proto/v1"
	"github.com/KimNattanan/go-chat-backend/internal/chat/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	v1.RoomServiceServer
	v1.MembershipServiceServer

	roomUseCase       usecase.RoomUseCase
	membershipUseCase usecase.MembershipUseCase
	l                 logger.Interface
	v                 *validator.Validate
}
