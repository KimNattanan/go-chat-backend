package rabbitmq

type Option func(*Server)

func Exchange(name string) Option {
	return func(s *Server) {
		s.exchange = name
	}
}

func RetryAttempts(attempts int) Option {
	return func(s *Server) {
		s.retryAttempts = attempts
	}
}
