package v1

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/chat/proto/v1"
	"github.com/KimNattanan/go-chat-backend/internal/chat/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	pbgrpc "google.golang.org/grpc"
)

// NewChatRoutes -.
func NewChatRoutes(app *pbgrpc.Server, roomUseCase usecase.RoomUseCase, membershipUseCase usecase.MembershipUseCase, l logger.Interface) {
	r := &V1{
		roomUseCase:       roomUseCase,
		membershipUseCase: membershipUseCase,
		l:                 l,
		v:                 validator.New(validator.WithRequiredStructEnabled()),
	}
	{
		v1.RegisterRoomServiceServer(app, r)
		v1.RegisterMembershipServiceServer(app, r)
	}
}
