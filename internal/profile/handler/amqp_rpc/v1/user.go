package v1

import (
	"context"
	"encoding/json"

	"github.com/KimNattanan/go-chat-backend/internal/profile/entity"
	"github.com/KimNattanan/go-chat-backend/internal/profile/handler/amqp_rpc/v1/request"
	"github.com/KimNattanan/go-chat-backend/pkg/responses"
	"github.com/google/uuid"
)

func (r *V1) userCreated(ctx context.Context, data []byte) error {
	var req request.UserCreatedRequest
	if err := json.Unmarshal(data, &req); err != nil {
		responses.LogAMQP(r.l, err, "amqp_rpc - V1 - userCreated")
		return err
	}
	if err := r.v.Struct(&req); err != nil {
		responses.LogAMQP(r.l, err, "amqp_rpc - V1 - userCreated")
		return err
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		responses.LogAMQP(r.l, err, "amqp_rpc - V1 - userCreated")
		return err
	}

	profile := &entity.Profile{
		UserID: userID,
		Email:  req.Email,
		Name:   req.Name,
	}
	if err := r.profileUsecase.Create(ctx, profile); err != nil {
		responses.LogAMQP(r.l, err, "amqp_rpc - V1 - userCreated")
		return err
	}

	return nil
}

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

	if err := r.profileUsecase.Delete(ctx, req.UserID); err != nil {
		responses.LogAMQP(r.l, err, "amqp_rpc - V1 - userDeleted")
		return err
	}

	return nil
}
