package request

type RoomDeletedRequest struct {
	RoomID string `json:"room_id" validate:"required,uuid"`
}
