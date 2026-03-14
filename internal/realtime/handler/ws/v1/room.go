package v1

import (
	"encoding/json"
	"errors"

	"github.com/KimNattanan/go-chat-backend/internal/platform/wsserver"
	"github.com/KimNattanan/go-chat-backend/internal/realtime/handler/ws/v1/request"
	"github.com/KimNattanan/go-chat-backend/pkg/responses"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

func (r *V1) roomWebSocket(c *echo.Context) error {
	roomID := c.Param("roomID")

	conn, err := r.wsServer.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		r.l.Error(err, "v1 - roomWebSocket")
		return responses.ErrorResponse(c, err)
	}
	r.wsServer.Register(roomID, conn)
	defer func() {
		r.wsServer.Unregister(roomID, conn)
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		msg, err := wsserver.ParseMessage(message)
		if err != nil {
			continue
		}
		switch msg.Type {
		case "create_message":
			var req request.CreateMessageRequest
			if err := json.Unmarshal(msg.Data, &req); err != nil {
				r.l.Error(err, "v1 - roomWebSocket")
				conn.WriteMessage(websocket.TextMessage, []byte("invalid request"))
				continue
			}
			if err := r.v.Struct(&req); err != nil {
				r.l.Error(err, "v1 - roomWebSocket")
				conn.WriteMessage(websocket.TextMessage, []byte("invalid request"))
				continue
			}
			messageID := uuid.New().String()
			r.amqpPublisher.Publish("message.created", map[string]string{
				"message_id": messageID,
				"room_id":    roomID,
				"user_id":    req.UserID,
				"content":    req.Content,
			})

		case "delete_message":
			var req request.DeleteMessageRequest
			if err := json.Unmarshal(msg.Data, &req); err != nil {
				r.l.Error(err, "v1 - roomWebSocket")
				conn.WriteMessage(websocket.TextMessage, []byte("invalid request"))
				continue
			}
			if err := r.v.Struct(&req); err != nil {
				r.l.Error(err, "v1 - roomWebSocket")
				conn.WriteMessage(websocket.TextMessage, []byte("invalid request"))
				continue
			}
			r.amqpPublisher.Publish("message.deleted", map[string]string{
				"room_id":    roomID,
				"message_id": req.MessageID,
			})

		default:
			r.l.Error(errors.New("invalid message type"), "v1 - roomWebSocket")
			conn.WriteMessage(websocket.TextMessage, []byte("invalid message type"))
			continue
		}
	}

	return nil
}
