package rabbitmq

import (
	"maps"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

const (
	_defaultExchangeName   = "app.fanout"
	_defaultReconnectDelay = 5 * time.Second
	_defaultRetryAttempts  = 3
)

type Handler func(ctx context.Context, data []byte) error

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc
	eg     *errgroup.Group

	url           string
	exchange      string
	retryAttempts int

	mu        sync.RWMutex
	conn      *amqp.Connection
	channel   *amqp.Channel
	consumers []consumerConfig

	notify chan error

	logger logger.Interface
}

type consumerConfig struct {
	queue   string
	workers int
	router  map[string]Handler
}

// New creates new Server instance.
func New(l logger.Interface, url string, opts ...Option) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(ctx)

	s := &Server{
		ctx:           ctx,
		cancel:        cancel,
		eg:            group,
		url:           url,
		exchange:      _defaultExchangeName,
		notify:        make(chan error, 1),
		logger:        l,
		retryAttempts: _defaultRetryAttempts,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) RegisterConsumer(queue string, workers int, router map[string]Handler) {
	if workers <= 0 {
		workers = 1
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.consumers = append(s.consumers, consumerConfig{
		queue:   queue,
		workers: workers,
		router:  router,
	})
}

func (s *Server) Start() {
	s.eg.Go(func() error {
		for {
			if err := s.runOnce(); err != nil {
				if errors.Is(err, context.Canceled) {
					return err
				}

				s.logger.Error(err, "rabbitmq - Server - runOnce")

				select {
				case s.notify <- err:
				default:
				}

				select {
				case <-time.After(_defaultReconnectDelay):
					continue
				case <-s.ctx.Done():
					return s.ctx.Err()
				}
			}

			return nil
		}
	})

	s.logger.Info("rabbitmq - Server - Started")
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	var shutdownErrors []error

	s.cancel()

	s.mu.Lock()
	if s.channel != nil {
		if err := s.channel.Close(); err != nil && !errors.Is(err, amqp.ErrClosed) {
			shutdownErrors = append(shutdownErrors, err)
		}
	}
	if s.conn != nil {
		if err := s.conn.Close(); err != nil && !errors.Is(err, amqp.ErrClosed) {
			shutdownErrors = append(shutdownErrors, err)
		}
	}
	s.mu.Unlock()

	if err := s.eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		shutdownErrors = append(shutdownErrors, err)
	}

	s.logger.Info("rabbitmq - Server - Shutdown")

	return errors.Join(shutdownErrors...)
}

func (s *Server) runOnce() error {
	conn, err := amqp.Dial(s.url)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}

	if err := ch.ExchangeDeclare(
		s.exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return err
	}

	if err := ch.ExchangeDeclare(
		"app.dlx",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return err
	}

	_, err = ch.QueueDeclare(
		"app.dlq",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return err
	}

	if err := ch.QueueBind(
		"app.dlq",
		"",
		"app.dlx",
		false,
		nil,
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return err
	}

	s.mu.Lock()
	s.conn = conn
	s.channel = ch
	consumers := append([]consumerConfig(nil), s.consumers...)
	s.mu.Unlock()

	for _, c := range consumers {
		if err := s.startConsumer(ch, c); err != nil {
			return err
		}
	}

	errCh := make(chan *amqp.Error, 1)
	ch.NotifyClose(errCh)

	select {
	case <-s.ctx.Done():
		return s.ctx.Err()
	case amqpErr := <-errCh:
		if amqpErr == nil {
			return nil
		}
		return amqpErr
	}
}

func (s *Server) startConsumer(ch *amqp.Channel, cfg consumerConfig) error {
	q, err := ch.QueueDeclare(
		cfg.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	if err := ch.QueueBind(
		q.Name,
		"",
		s.exchange,
		false,
		nil,
	); err != nil {
		return err
	}

	deliveries, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for i := 0; i < cfg.workers; i++ {
		go s.worker(deliveries, cfg.router)
	}

	return nil
}

func (s *Server) worker(
	deliveries <-chan amqp.Delivery,
	router map[string]Handler,
) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case d, ok := <-deliveries:
			if !ok {
				return
			}

			msg, err := ParseMessage(d.Body)
			if err != nil {
				d.Nack(false, false)
				continue
			}

			handler, ok := router[msg.Type]
			if !ok {
				d.Nack(false, false)
				continue
			}

			ctx, cancel := context.WithTimeout(s.ctx, 5*time.Minute)
			err = handler(ctx, msg.Data)
			cancel()
			if err != nil {
				if retryErr := s.retry(&d); retryErr != nil {
					s.logger.Error(retryErr, "rabbitmq - worker - retry failed")
					d.Nack(false, false)
				}
				continue
			}

			d.Ack(false)
		}
	}
}

func getRetryCount(headers amqp.Table) int {
	v, ok := headers["x-retry-count"]
	if !ok {
		return 0
	}
	n, ok := v.(int)
	if !ok {
		return 0
	}
	return n
}

func (s *Server) retry(d *amqp.Delivery) error {
	retryCount := getRetryCount(d.Headers) + 1

	// send to DLQ
	if retryCount >= s.retryAttempts {
		dlqHeaders := amqp.Table{}
		maps.Copy(dlqHeaders, d.Headers)
		dlqHeaders["x-retry-count"] = retryCount
		if err := s.channel.Publish(
			"app.dlx",
			"",
			false,
			false,
			amqp.Publishing{
				Headers:     dlqHeaders,
				ContentType: d.ContentType,
				Body:        d.Body,
			},
		); err != nil {
			d.Nack(false, false)
			return err
		}
		d.Ack(false)
		return nil
	}

	// republish

	headers := amqp.Table{}
	maps.Copy(headers, d.Headers)
	headers["x-retry-count"] = retryCount

	exchange := d.Exchange
	if exchange == "" {
		exchange = s.exchange
	}

	if err := s.channel.Publish(
		exchange,
		d.RoutingKey,
		false,
		false,
		amqp.Publishing{
			Headers:     headers,
			ContentType: d.ContentType,
			Body:        d.Body,
		},
	); err != nil {
		d.Nack(false, false)
		return err
	}
	d.Ack(false)
	return nil
}
