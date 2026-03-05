package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/streadway/amqp"
)

const (
	workerTimeout = 5 * 60 * time.Second
)

type Handler func(ctx context.Context, msgType string, data []byte) error

func (c *Client) Consume(workers int, handler Handler) error {
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

	jobs := make(chan amqp.Delivery, 5)
	defer close(jobs)

	for range workers {
		go c.worker(jobs, handler)
	}

	for {
		d, ok := <-deliveries
		if !ok { // deliveries channel closed
			return nil
		}
		jobs <- d
	}
}

func (c *Client) worker(
	jobs <-chan amqp.Delivery,
	handler Handler,
) {
	for d := range jobs {
		var msg Message
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			d.Nack(false, false)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), workerTimeout)
		err := handler(ctx, msg.Type, msg.Data)
		cancel()

		if err != nil {
			d.Nack(false, true)
			continue
		}

		d.Ack(false)
	}
}
