package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/chat/entity"
	v1 "github.com/KimNattanan/go-chat-backend/internal/chat/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoRoom(o *entity.Room) *v1.Room {
	if o == nil {
		return nil
	}
	memberships := make([]*v1.RoomMembership, len(o.Memberships))
	for i := range o.Memberships {
		memberships[i].UserId = o.Memberships[i].UserID.String()
		memberships[i].CreatedAt = timestamppb.New(o.Memberships[i].CreatedAt)
	}
	return &v1.Room{
		Id:        o.ID.String(),
		Title:     o.Title,
		CreatedAt: timestamppb.New(o.CreatedAt),
		UpdatedAt: timestamppb.New(o.UpdatedAt),

		Memberships: memberships,
	}
}

func toProtoRoomList(rooms []*entity.Room) []*v1.Room {
	result := make([]*v1.Room, 0, len(rooms))
	for _, o := range rooms {
		result = append(result, toProtoRoom(o))
	}
	return result
}

func toProtoMembership(o *entity.Membership) *v1.Membership {
	if o == nil {
		return nil
	}
	return &v1.Membership{
		RoomId:    o.RoomID.String(),
		UserId:    o.UserID.String(),
		CreatedAt: timestamppb.New(o.CreatedAt),
	}
}

func toProtoMembershipList(memberships []*entity.Membership) []*v1.Membership {
	result := make([]*v1.Membership, 0, len(memberships))
	for _, o := range memberships {
		result = append(result, toProtoMembership(o))
	}
	return result
}
