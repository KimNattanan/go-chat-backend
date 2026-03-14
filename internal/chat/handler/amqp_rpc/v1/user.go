package v1

import (
	"context"
	"encoding/json"

	"github.com/KimNattanan/go-chat-backend/internal/chat/handler/amqp_rpc/v1/request"
)

func (r *V1) userDeleted(ctx context.Context, data []byte) error {
	var req request.UserDeletedRequest
	if err := json.Unmarshal(data, &req); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - userDeleted")
		return err
	}
	if err := r.v.Struct(&req); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - userDeleted")
		return err
	}

	if err := r.membershipUseCase.DeleteByUserID(ctx, req.UserID); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - userDeleted")
		return err
	}

	return nil
}
