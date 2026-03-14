package request

type UserCreatedRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Email  string `json:"email" validate:"required,email"`
	Name   string `json:"name"`
}

type UserDeletedRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}
