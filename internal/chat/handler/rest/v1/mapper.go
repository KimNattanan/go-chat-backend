package v1

import (
	"time"

	"github.com/KimNattanan/go-chat-backend/internal/chat/entity"
	"github.com/KimNattanan/go-chat-backend/internal/chat/handler/rest/v1/response"
	"github.com/google/uuid"
)

func toRoomResponse(r *entity.Room) *response.RoomResponse {
	memberships := make(
		[]*struct {
			UserID    uuid.UUID `json:"user_id"`
			CreatedAt time.Time `json:"created_at"`
		},
		len(r.Memberships),
	)
	for i := range r.Memberships {
		memberships[i].UserID = r.Memberships[i].UserID
		memberships[i].CreatedAt = r.Memberships[i].CreatedAt
	}

	return &response.RoomResponse{
		ID:          r.ID,
		Title:       r.Title,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		Memberships: memberships,
	}
}

func toRoomResponseList(rooms []*entity.Room) []*response.RoomResponse {
	result := make([]*response.RoomResponse, 0, len(rooms))
	for _, o := range rooms {
		result = append(result, toRoomResponse(o))
	}
	return result
}

func toMembershipResponse(membership *entity.Membership) *response.MembershipResponse {
	return &response.MembershipResponse{
		RoomID:    membership.RoomID,
		UserID:    membership.UserID,
		CreatedAt: membership.CreatedAt,
	}
}

func toMembershipResponseList(memberships []*entity.Membership) []*response.MembershipResponse {
	result := make([]*response.MembershipResponse, 0, len(memberships))
	for _, o := range memberships {
		result = append(result, toMembershipResponse(o))
	}
	return result
}
