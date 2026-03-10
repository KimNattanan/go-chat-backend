package v1

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/message/entity"
	v1 "github.com/KimNattanan/go-chat-backend/internal/message/proto/v1"
	"github.com/KimNattanan/go-chat-backend/pkg/apperror"
	"github.com/google/uuid"
)

func (r *V1) FindMessageByID(ctx context.Context, req *v1.FindMessageByIDRequest) (*v1.MessageResponse, error) {
	message, err := r.messageUseCase.FindByID(ctx, req.Id)
	if err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.MessageResponse{
		Message: toProtoMessage(message),
	}, nil
}

func (r *V1) FindMessageByRoomID(ctx context.Context, req *v1.FindMessagesByRoomIDRequest) (*v1.MessagesResponse, error) {
	messages, err := r.messageUseCase.FindByRoomID(ctx, req.RoomId)
	if err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.MessagesResponse{
		Messages: toProtoMessageList(messages),
	}, nil
}

func (r *V1) FindMessageByUserID(ctx context.Context, req *v1.FindMessagesByUserIDRequest) (*v1.MessagesResponse, error) {
	messages, err := r.messageUseCase.FindByUserID(ctx, req.UserId)
	if err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.MessagesResponse{
		Messages: toProtoMessageList(messages),
	}, nil
}

func (r *V1) FindMessageByRoomIDAndUserID(ctx context.Context, req *v1.FindMessagesByRoomIDAndUserIDRequest) (*v1.MessagesResponse, error) {
	messages, err := r.messageUseCase.FindByRoomIDAndUserID(ctx, req.RoomId, req.UserId)
	if err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.MessagesResponse{
		Messages: toProtoMessageList(messages),
	}, nil
}

func (r *V1) CreateMessage(ctx context.Context, req *v1.CreateMessageRequest) (*v1.MessageResponse, error) {
	roomID, err := uuid.Parse(req.RoomId)
	if err != nil {
		return nil, apperror.ParseGrpc(err)
	}
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	message := &entity.Message{
		RoomID:  roomID,
		UserID:  userID,
		Content: req.Content,
	}
	if err := r.messageUseCase.Create(ctx, message); err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.MessageResponse{
		Message: toProtoMessage(message),
	}, nil
}

func (r *V1) DeleteMessage(ctx context.Context, req *v1.DeleteMessageRequest) (*v1.DeleteMessageResponse, error) {
	if err := r.messageUseCase.Delete(ctx, req.Id); err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.DeleteMessageResponse{
		Message: "message deleted",
	}, nil
}

func (r *V1) DeleteMessagesByRoomID(ctx context.Context, req *v1.DeleteMessagesByRoomIDRequest) (*v1.DeleteMessageResponse, error) {
	if err := r.messageUseCase.DeleteByRoomID(ctx, req.RoomId); err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.DeleteMessageResponse{
		Message: "messages deleted",
	}, nil
}

func (r *V1) AnonymizeUserMessages(ctx context.Context, req *v1.AnonymizeUserMessagesRequest) (*v1.AnonymizeUserMessagesResponse, error) {
	if err := r.messageUseCase.AnonymizeUserMessages(ctx, req.UserId); err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.AnonymizeUserMessagesResponse{
		Message: "messages anonymized",
	}, nil
}
