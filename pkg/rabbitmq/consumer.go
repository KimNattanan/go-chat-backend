package rabbitmq

import (
	"context"
	"time"

	"github.com/streadway/amqp"
)

const (
	workerTimeout = 5 * 60 * time.Second
)

type Handler func(ctx context.Context, d *amqp.Delivery, data []byte)

func (c *Client) Consume(workers int, router map[string]Handler) error {
	c.mu.RLock()
	ch := c.consCh
	queue := c.queue
	c.mu.RUnlock()

	deliveries, err := ch.Consume(
		queue,
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

	for range workers {
		go c.worker(deliveries, router)
	}
	return nil
}

func (c *Client) worker(
	deliveries <-chan amqp.Delivery,
	router map[string]Handler,
) {
	for d := range deliveries {
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

		ctx, cancel := context.WithTimeout(context.Background(), workerTimeout)
		handler(ctx, &d, msg.Data)
		cancel()
	}
}
