package v1

import (
	"context"
	"encoding/json"

	"github.com/KimNattanan/go-chat-backend/internal/profile/entity"
	"github.com/KimNattanan/go-chat-backend/internal/profile/handler/amqp_rpc/v1/request"
	"github.com/google/uuid"
)

func (r *V1) createProfile(ctx context.Context, data []byte) error {
	var req request.CreateProfileRequest
	if err := json.Unmarshal(data, &req); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - createProfile")
		return err
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		r.l.Error(err, "amqp_rpc - V1 - createProfile")
		return err
	}

	profile := &entity.Profile{
		UserID: userID,
		Email:  req.Email,
		Name:   req.Name,
	}
	if err := r.profileUsecase.Create(ctx, profile); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - createProfile")
		return err
	}

	return nil
}

func (r *V1) deleteProfile(ctx context.Context, data []byte) error {
	var req request.DeleteProfileRequest
	if err := json.Unmarshal(data, &req); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - createProfile")
		return err
	}

	if err := r.profileUsecase.Delete(ctx, req.UserID); err != nil {
		r.l.Error(err, "amqp_rpc - V1 - createProfile")
		return err
	}

	return nil
}
