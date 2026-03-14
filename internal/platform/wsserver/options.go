package wsserver

import "net/http"

type Option func(*Server)

func CheckOrigin(fn func(r *http.Request) bool) Option {
	return func(s *Server) {
		s.Upgrader.CheckOrigin = fn
	}
}
