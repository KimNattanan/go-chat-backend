package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/auth/entity"
	"github.com/KimNattanan/go-chat-backend/internal/auth/handler/rest/v1/response"
)

func toUserResponse(p *entity.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        p.ID,
		Email:     p.Email,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
