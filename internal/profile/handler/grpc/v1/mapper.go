package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/profile/entity"
	v1 "github.com/KimNattanan/go-chat-backend/internal/profile/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoProfile(o *entity.Profile) *v1.Profile {
	if o == nil {
		return nil
	}
	return &v1.Profile{
		UserId:    o.UserID.String(),
		Email:     o.Email,
		Name:      o.Name,
		CreatedAt: timestamppb.New(o.CreatedAt),
		UpdatedAt: timestamppb.New(o.UpdatedAt),
	}
}
