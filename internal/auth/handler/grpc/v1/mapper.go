package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/auth/entity"
	v1 "github.com/KimNattanan/go-chat-backend/internal/auth/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoUser(o *entity.User) *v1.User {
	if o == nil {
		return nil
	}
	return &v1.User{
		Id:        o.ID.String(),
		Email:     o.Email,
		CreatedAt: timestamppb.New(o.CreatedAt),
		UpdatedAt: timestamppb.New(o.UpdatedAt),
	}
}
