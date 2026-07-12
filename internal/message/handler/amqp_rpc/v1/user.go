package v1

import (
	"context"
	"encoding/json"

	"github.com/KimNattanan/go-chat-backend/internal/message/handler/amqp_rpc/v1/request"
	"github.com/KimNattanan/go-chat-backend/pkg/responses"
)

func (r *V1) userDeleted(ctx context.Context, data []byte) error {
	var req request.UserDeletedRequest
	if err := json.Unmarshal(data, &req); err != nil {
		responses.LogAMQP(r.l, err, "amqp_rpc - V1 - userDeleted")
		return err
	}
	if err := r.v.Struct(&req); err != nil {
		responses.LogAMQP(r.l, err, "amqp_rpc - V1 - userDeleted")
		return err
	}

	if err := r.messageUseCase.AnonymizeUserMessages(ctx, req.UserID); err != nil {
		responses.LogAMQP(r.l, err, "amqp_rpc - V1 - userDeleted")
		return err
	}

	return nil
}
