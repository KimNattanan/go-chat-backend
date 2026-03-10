package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/message/entity"
	v1 "github.com/KimNattanan/go-chat-backend/internal/message/proto/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoMessage(o *entity.Message) *v1.Message {
	if o == nil {
		return nil
	}
	return &v1.Message{
		Id:        o.ID.String(),
		RoomId:    o.RoomID.String(),
		UserId:    o.UserID.String(),
		Content:   o.Content,
		CreatedAt: timestamppb.New(o.CreatedAt),
	}
}

func toProtoMessageList(messages []*entity.Message) []*v1.Message {
	result := make([]*v1.Message, 0, len(messages))
	for _, o := range messages {
		result = append(result, toProtoMessage(o))
	}
	return result
}
