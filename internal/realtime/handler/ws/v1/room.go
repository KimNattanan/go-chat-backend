package v1

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/KimNattanan/go-chat-backend/internal/platform/wsserver"
	"github.com/KimNattanan/go-chat-backend/internal/realtime/handler/ws/v1/request"
	"github.com/KimNattanan/go-chat-backend/pkg/responses"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

const (
	pongWait   = 60 * time.Second
	writeWait  = 10 * time.Second
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

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPingHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(writeWait))
	})

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.SetReadDeadline(time.Now().Add(pongWait))

		// Ignore control frames
		if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
			continue
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
