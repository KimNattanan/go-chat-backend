package rest

import (
	v1 "github.com/KimNattanan/go-chat-backend/internal/chat/handler/grpc/v1"
	"github.com/KimNattanan/go-chat-backend/internal/chat/usecase"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	pbgrpc "google.golang.org/grpc"
)

// NewRouter -.
func NewRouter(app *pbgrpc.Server, roomUseCase usecase.RoomUseCase, membershipUseCase usecase.MembershipUseCase, l logger.Interface) {
	v1.NewChatRoutes(app, roomUseCase, membershipUseCase, l)
}
