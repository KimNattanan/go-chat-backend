package request

type CreateProfileRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required"`
}

type PatchProfileRequest struct {
	Name string `json:"name" validate:"required"`
}
