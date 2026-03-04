package v1

import (
	"github.com/KimNattanan/go-chat-backend/internal/profile/entity"
	"github.com/KimNattanan/go-chat-backend/internal/profile/handler/rest/v1/response"
)

func toProfileResponse(p *entity.Profile) *response.ProfileResponse {
	return &response.ProfileResponse{
		UserID:    p.UserID,
		Email:     p.Email,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
