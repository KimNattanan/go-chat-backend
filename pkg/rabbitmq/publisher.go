package rabbitmq

import (
	"errors"

	"github.com/streadway/amqp"
)

func (c *Client) Publish(msgType string, data any) error {
	c.mu.RLock()
	ch := c.pubCh
	queue := c.queue
	c.mu.RUnlock()

	if ch == nil {
		return errors.New("publisher channel not ready")
	}

	body, err := NewMessage(msgType, data)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
