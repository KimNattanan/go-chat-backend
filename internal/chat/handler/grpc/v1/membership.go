package v1

import (
	"context"

	"github.com/KimNattanan/go-chat-backend/internal/chat/entity"
	v1 "github.com/KimNattanan/go-chat-backend/internal/chat/proto/v1"
	"github.com/KimNattanan/go-chat-backend/pkg/apperror"
	"github.com/google/uuid"
)

func (r *V1) FindMembershipsByRoomID(ctx context.Context, req *v1.FindMembershipsByRoomIDRequest) (*v1.MembershipsResponse, error) {
	memberships, err := r.membershipUseCase.FindByRoomID(ctx, req.RoomId)
	if err != nil {
		return nil, apperror.ParseGrpcLogged(ctx, r.l, err, "grpc - v1 - findMembershipsByRoomID")
	}

	return &v1.MembershipsResponse{
		Memberships: toProtoMembershipList(memberships),
	}, nil
}

func (r *V1) FindMembershipsByUserID(ctx context.Context, req *v1.FindMembershipsByUserIDRequest) (*v1.MembershipsResponse, error) {
	memberships, err := r.membershipUseCase.FindByRoomID(ctx, req.UserId)
	if err != nil {
		return nil, apperror.ParseGrpcLogged(ctx, r.l, err, "grpc - v1 - findMembershipsByUserID")
	}

	return &v1.MembershipsResponse{
		Memberships: toProtoMembershipList(memberships),
	}, nil
}

func (r *V1) FindMembershipByRoomIDAndUserID(ctx context.Context, req *v1.FindMembershipByRoomIDAndUserIDRequest) (*v1.MembershipResponse, error) {
	membership, err := r.membershipUseCase.FindByRoomIDAndUserID(ctx, req.RoomId, req.UserId)
	if err != nil {
		return nil, apperror.ParseGrpcLogged(ctx, r.l, err, "grpc - v1 - findMembershipByRoomIDAndUserID")
	}

	return &v1.MembershipResponse{
		Membership: toProtoMembership(membership),
	}, nil
}

func (r *V1) CreateMembership(ctx context.Context, req *v1.CreateMembershipRequest) (*v1.MembershipResponse, error) {
	roomID, err := uuid.Parse(req.RoomId)
	if err != nil {
		return nil, apperror.ParseGrpcLogged(ctx, r.l, err, "grpc - v1 - createMembership")
	}
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, apperror.ParseGrpcLogged(ctx, r.l, err, "grpc - v1 - createMembership")
	}
	membership := &entity.Membership{
		RoomID: roomID,
		UserID: userID,
	}
	if err := r.membershipUseCase.Create(ctx, membership); err != nil {
		return nil, apperror.ParseGrpcLogged(ctx, r.l, err, "grpc - v1 - createMembership")
	}

	return &v1.MembershipResponse{
		Membership: toProtoMembership(membership),
	}, nil
}

func (r *V1) DeleteMembershipsByUserID(ctx context.Context, req *v1.DeleteMembershipsByUserIDRequest) (*v1.DeleteMembershipResponse, error) {
	if err := r.membershipUseCase.DeleteByUserID(ctx, req.UserId); err != nil {
		return nil, apperror.ParseGrpcLogged(ctx, r.l, err, "grpc - v1 - deleteMembershipsByUserID")
	}

	return &v1.DeleteMembershipResponse{
		Message: "memberships deleted",
	}, nil
}

func (r *V1) DeleteMembershipByRoomIDAndUserID(ctx context.Context, req *v1.DeleteMembershipByRoomIDAndUserIDRequest) (*v1.DeleteMembershipResponse, error) {
	if err := r.membershipUseCase.Delete(ctx, req.RoomId, req.UserId); err != nil {
		return nil, apperror.ParseGrpcLogged(ctx, r.l, err, "grpc - v1 - deleteMembershipByRoomIDAndUserID")
	}

	return &v1.DeleteMembershipResponse{
		Message: "membership deleted",
	}, nil
}
