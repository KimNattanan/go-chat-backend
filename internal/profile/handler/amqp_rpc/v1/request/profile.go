package request

type CreateProfileRequest struct {
	UserID string `json:"user_id"`
	Email  string `json:"email" validate:"required,email"`
	Name   string `json:"name"`
}

type DeleteProfileRequest struct {
	UserID string `json:"user_id"`
}
