package v1

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/chat/entity"
	v1 "github.com/KimNattanan/go-chat-backend/internal/chat/proto/v1"
	"github.com/KimNattanan/go-chat-backend/pkg/apperror"
)

func (r *V1) FindRoomByID(ctx context.Context, req *v1.FindRoomByIDRequest) (*v1.RoomResponse, error) {
	room, err := r.roomUseCase.FindByID(ctx, req.Id)
	if err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.RoomResponse{
		Room: toProtoRoom(room),
	}, nil
}

func (r *V1) FindRoomsByUserID(ctx context.Context, req *v1.FindRoomsByUserIDRequest) (*v1.RoomsResponse, error) {
	rooms, err := r.roomUseCase.FindByUserID(ctx, req.UserId)
	if err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.RoomsResponse{
		Rooms: toProtoRoomList(rooms),
	}, nil
}

func (r *V1) CreateRoom(ctx context.Context, req *v1.CreateRoomRequest) (*v1.RoomResponse, error) {
	room := &entity.Room{
		Title: req.Title,
	}
	if err := r.roomUseCase.Create(ctx, room); err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.RoomResponse{
		Room: toProtoRoom(room),
	}, nil
}

func (r *V1) PatchRoom(ctx context.Context, req *v1.PatchRoomRequest) (*v1.RoomResponse, error) {
	room := &entity.Room{
		Title: req.Title,
	}
	updatedRoom, err := r.roomUseCase.Patch(ctx, req.Id, room)
	if err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.RoomResponse{
		Room: toProtoRoom(updatedRoom),
	}, nil
}

func (r *V1) DeleteRoom(ctx context.Context, req *v1.DeleteRoomRequest) (*v1.DeleteRoomResponse, error) {
	if err := r.roomUseCase.Delete(ctx, req.Id); err != nil {
		return nil, apperror.ParseGrpc(err)
	}

	return &v1.DeleteRoomResponse{
		Message: "room deleted",
	}, nil
}
