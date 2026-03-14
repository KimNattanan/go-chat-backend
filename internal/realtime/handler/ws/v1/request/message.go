package request

type CreateMessageRequest struct {
	UserID  string `json:"user_id" validate:"required,uuid"`
	Content string `json:"content" validate:"required"`
}

type DeleteMessageRequest struct {
	MessageID string `json:"message_id" validate:"required,uuid"`
}
