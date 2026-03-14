package v1

import (
	"context"
	"encoding/json"

	"github.com/KimNattanan/go-chat-backend/internal/message/handler/amqp_rpc/v1/request"
)

func (r *V1) messageCreated(ctx context.Context, data []byte) error {
	var req request.MessageCreatedRequest
	if err := json.Unmarshal(data, &req); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - messageCreated")
		return err
	}
	if err := r.v.Struct(&req); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - messageCreated")
		return err
	}

	r.wsServer.BroadcastMessage(req.RoomID, "create_message", map[string]string{
		"user_id": req.UserID,
		"content": req.Content,
	})

	return nil
}

func (r *V1) messageDeleted(ctx context.Context, data []byte) error {
	var req request.MessageDeletedRequest
	if err := json.Unmarshal(data, &req); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - messageDeleted")
		return err
	}
	if err := r.v.Struct(&req); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - messageDeleted")
		return err
	}

	r.wsServer.BroadcastMessage(req.RoomID, "delete_message", map[string]string{
		"message_id": req.MessageID,
	})

	return nil
}
