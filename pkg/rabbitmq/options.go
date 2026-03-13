package rabbitmq

// Option configures Server.
type Option func(*Server)

// Exchange sets the exchange name.
func Exchange(name string) Option {
	return func(s *Server) {
		s.exchange = name
	}
}
