package wsserver

import (
	"net/http"

	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/KimNattanan/go-chat-backend/pkg/rabbitmq"
	"github.com/gorilla/websocket"
)

type Server struct {
	logger        logger.Interface
	Upgrader      *websocket.Upgrader
	clients       map[string]map[*websocket.Conn]bool
	amqpPublisher *rabbitmq.Publisher
}

func New(l logger.Interface, opts ...Option) *Server {
	s := &Server{
		logger: l,
		Upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients: make(map[string]map[*websocket.Conn]bool),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) BroadcastMessage(roomID, msgType string, data any) {
	for client := range s.clients[roomID] {
		body, err := NewMessage(msgType, data)
		if err != nil {
			continue
		}
		if err := client.WriteMessage(websocket.TextMessage, body); err != nil {
			continue
		}
	}
}
