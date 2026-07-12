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
	pongWait  = 60 * time.Second
	writeWait = 10 * time.Second
)

func (r *V1) roomWebSocket(c *echo.Context) error {
	roomID := c.Param("roomID")

	conn, err := r.wsServer.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return responses.LogAndErrorResponse(c, r.l, err, "v1 - roomWebSocket")
	}
	r.wsServer.Register(roomID, conn)
	defer func() {
		r.wsServer.Unregister(roomID, conn)
		_ = conn.Close()
	}()

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPingHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(writeWait))
	})

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				r.l.Warn(err, "v1 - roomWebSocket - ReadMessage")
			}
			break
		}
		conn.SetReadDeadline(time.Now().Add(pongWait))

		// Ignore control frames
		if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
			continue
		}

		msg, err := wsserver.ParseMessage(message)
		if err != nil {
			r.l.Warn(err, "v1 - roomWebSocket - ParseMessage")
			continue
		}
		switch msg.Type {
		case "create_message":
			var req request.CreateMessageRequest
			if err := json.Unmarshal(msg.Data, &req); err != nil {
				r.l.Warn(err, "v1 - roomWebSocket")
				if werr := conn.WriteMessage(websocket.TextMessage, []byte("invalid request")); werr != nil {
					r.l.Error(werr, "v1 - roomWebSocket - WriteMessage")
				}
				continue
			}
			if err := r.v.Struct(&req); err != nil {
				r.l.Warn(err, "v1 - roomWebSocket")
				if werr := conn.WriteMessage(websocket.TextMessage, []byte("invalid request")); werr != nil {
					r.l.Error(werr, "v1 - roomWebSocket - WriteMessage")
				}
				continue
			}
			messageID := uuid.New().String()
			if err := r.mqPublisher.Publish("message.created", map[string]string{
				"message_id": messageID,
				"room_id":    roomID,
				"user_id":    req.UserID,
				"content":    req.Content,
			}); err != nil {
				r.l.Error(err, "v1 - roomWebSocket - Publish create_message")
			}

		case "delete_message":
			var req request.DeleteMessageRequest
			if err := json.Unmarshal(msg.Data, &req); err != nil {
				r.l.Warn(err, "v1 - roomWebSocket")
				if werr := conn.WriteMessage(websocket.TextMessage, []byte("invalid request")); werr != nil {
					r.l.Error(werr, "v1 - roomWebSocket - WriteMessage")
				}
				continue
			}
			if err := r.v.Struct(&req); err != nil {
				r.l.Warn(err, "v1 - roomWebSocket")
				if werr := conn.WriteMessage(websocket.TextMessage, []byte("invalid request")); werr != nil {
					r.l.Error(werr, "v1 - roomWebSocket - WriteMessage")
				}
				continue
			}
			if err := r.mqPublisher.Publish("message.deleted", map[string]string{
				"room_id":    roomID,
				"message_id": req.MessageID,
			}); err != nil {
				r.l.Error(err, "v1 - roomWebSocket - Publish delete_message")
			}

		default:
			r.l.Warn(errors.New("invalid message type"), "v1 - roomWebSocket")
			if werr := conn.WriteMessage(websocket.TextMessage, []byte("invalid message type")); werr != nil {
				r.l.Error(werr, "v1 - roomWebSocket - WriteMessage")
			}
			continue
		}
	}

	return nil
}
