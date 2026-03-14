package request

type MessageCreatedRequest struct {
	MessageID string `json:"message_id" validate:"required,uuid"`
	RoomID    string `json:"room_id" validate:"required,uuid"`
	UserID    string `json:"user_id" validate:"required,uuid"`
	Content   string `json:"content" validate:"required"`
}

type MessageDeletedRequest struct {
	RoomID    string `json:"room_id" validate:"required,uuid"`
	MessageID string `json:"message_id" validate:"required,uuid"`
}
