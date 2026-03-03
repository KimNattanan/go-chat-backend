package request

type PatchProfileRequest struct {
	Name string `json:"name" validate:"required"`
}
