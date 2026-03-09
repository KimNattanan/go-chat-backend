package request

type CreateRoomRequest struct {
	Title string `json:"title"`
}

type PatchRoomRequest struct {
	Title string `json:"title"`
}
