package v1

import (
	"context"
	"encoding/json"

	"github.com/KimNattanan/go-chat-backend/internal/message/entity"
	"github.com/KimNattanan/go-chat-backend/internal/message/handler/amqp_rpc/v1/request"
	"github.com/google/uuid"
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

	messageID, err := uuid.Parse(req.MessageID)
	if err != nil {
		r.l.Error(err, "amqp_rpc - V1 - messageCreated")
		return err
	}
	roomID, err := uuid.Parse(req.RoomID)
	if err != nil {
		r.l.Error(err, "amqp_rpc - V1 - messageCreated")
		return err
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		r.l.Error(err, "amqp_rpc - V1 - messageCreated")
		return err
	}

	message := &entity.Message{
		ID:      messageID,
		RoomID:  roomID,
		UserID:  userID,
		Content: req.Content,
	}
	if err := r.messageUseCase.Create(ctx, message); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - messageCreated")
		return err
	}

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

	if err := r.messageUseCase.Delete(ctx, req.MessageID); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - messageDeleted")
		return err
	}

	return nil
}
