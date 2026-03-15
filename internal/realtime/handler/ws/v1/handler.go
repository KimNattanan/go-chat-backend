package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/platform/wsserver"
	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	wsServer    *wsserver.Server
	mqPublisher rabbitmq.Publisher

	l logger.Interface
	v *validator.Validate
}
