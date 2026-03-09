package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/message/entity"
	"github.com/KimNattanan/go-chat-backend/internal/message/handler/rest/v1/response"
)

func toMessageResponse(m *entity.Message) *response.MessageResponse {
	return &response.MessageResponse{
		ID:        m.ID,
		RoomID:    m.RoomID,
		UserID:    m.UserID,
		CreatedAt: m.CreatedAt,
	}
}

func toMessageResponseList(rooms []*entity.Message) []*response.MessageResponse {
	result := make([]*response.MessageResponse, 0, len(rooms))
	for _, o := range rooms {
		result = append(result, toMessageResponse(o))
	}
	return result
}
