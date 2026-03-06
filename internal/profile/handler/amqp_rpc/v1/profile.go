package v1

import (
	"context"
	"encoding/json"

	"github.com/KimNattanan/go-chat-backend/internal/profile/entity"
	"github.com/KimNattanan/go-chat-backend/internal/profile/handler/amqp_rpc/v1/request"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func (r *V1) createProfile() rabbitmq.Handler {
	return func(ctx context.Context, d *amqp.Delivery, data []byte) {
		var req request.CreateProfileRequest
		if err := json.Unmarshal(data, &req); err != nil {
			r.l.Error(err, "amqp_rpc - V1 - createProfile")
			d.Nack(false, false)
			return
		}
		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			r.l.Error(err, "amqp_rpc - V1 - createProfile")
			d.Nack(false, false)
			return
		}
		profile := &entity.Profile{
			UserID: userID,
			Email:  req.Email,
			Name:   req.Name,
		}
		if err := r.profileUsecase.Create(ctx, profile); err != nil {
			r.l.Error(err, "amqp_rpc - V1 - createProfile")
			d.Nack(false, false)
			return
		}
		d.Ack(false)
	}
}

func (r *V1) deleteProfile() rabbitmq.Handler {
	return func(ctx context.Context, d *amqp.Delivery, data []byte) {
		var req request.DeleteProfileRequest
		if err := json.Unmarshal(data, &req); err != nil {
			r.l.Error(err, "amqp_rpc - V1 - createProfile")
			d.Nack(false, false)
			return
		}
		if err := r.profileUsecase.Delete(ctx, req.UserID); err != nil {
			r.l.Error(err, "amqp_rpc - V1 - createProfile")
			d.Nack(false, false)
			return
		}
		d.Ack(false)
	}
}
