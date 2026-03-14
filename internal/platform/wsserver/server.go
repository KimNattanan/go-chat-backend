package wsserver

import (
	"net/http"
	"sync"

	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/gorilla/websocket"
)

type Server struct {
	logger       logger.Interface
	Upgrader     *websocket.Upgrader
	clients      map[string]map[*websocket.Conn]bool
	clientsMutex sync.RWMutex
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
	body, err := NewMessage(msgType, data)
	if err != nil {
		return
	}

	s.clientsMutex.RLock()
	roomClients := s.clients[roomID]
	if len(roomClients) == 0 {
		s.clientsMutex.RUnlock()
		return
	}
	conns := make([]*websocket.Conn, 0, len(roomClients))
	for conn := range roomClients {
		conns = append(conns, conn)
	}
	s.clientsMutex.RUnlock()

	for _, client := range conns {
		if err := client.WriteMessage(websocket.TextMessage, body); err != nil {
			s.Unregister(roomID, client)
		}
	}
}

func (s *Server) Register(roomID string, conn *websocket.Conn) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	if _, ok := s.clients[roomID]; !ok {
		s.clients[roomID] = make(map[*websocket.Conn]bool)
	}

	s.clients[roomID][conn] = true
}

func (s *Server) Unregister(roomID string, conn *websocket.Conn) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	if _, ok := s.clients[roomID]; !ok {
		return
	}
	if _, ok := s.clients[roomID][conn]; !ok {
		return
	}

	delete(s.clients[roomID], conn)
	if len(s.clients[roomID]) == 0 {
		delete(s.clients, roomID)
	}
}
