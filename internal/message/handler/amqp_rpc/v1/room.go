package v1

import (
	"context"
	"encoding/json"

	"github.com/KimNattanan/go-chat-backend/internal/message/handler/amqp_rpc/v1/request"
)

func (r *V1) roomDeleted(ctx context.Context, data []byte) error {
	var req request.RoomDeletedRequest
	if err := json.Unmarshal(data, &req); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - roomDeleted")
		return err
	}
	if err := r.v.Struct(&req); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - roomDeleted")
		return err
	}

	if err := r.messageUseCase.DeleteByRoomID(ctx, req.RoomID); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - roomDeleted")
		return err
	}

	return nil
}
