package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/KimNattanan/go-chat-backend/pkg/logger"
	"github.com/streadway/amqp"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type Client struct {
	url            string
	queue          string
	conn           *amqp.Connection
	channel        *amqp.Channel
	closeChan      chan *amqp.Error
	mutex          sync.Mutex
	reconnectDelay time.Duration
	logger         logger.Interface
}

// defer client.Close()
func New(url, queue string, reconnectDelay int, l logger.Interface) *Client {
	return &Client{
		url:            url,
		queue:          queue,
		reconnectDelay: time.Duration(reconnectDelay) * time.Second,
		logger:         l,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	for {
		conn, err := amqp.Dial(c.url)
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(c.reconnectDelay):
				continue
			}
		}

		ch, err := conn.Channel()
		if err != nil {
			conn.Close()
			continue
		}

		_, err = ch.QueueDeclare(
			c.queue,
			true,  // durable
			false, // auto delete
			false, // exclusive
			false, // no wait
			nil,
		)
		if err != nil {
			conn.Close()
			continue
		}

		c.mutex.Lock()
		c.conn = conn
		c.channel = ch
		c.closeChan = make(chan *amqp.Error)
		conn.NotifyClose(c.closeChan)
		c.mutex.Unlock()

		return nil
	}
}

func (c *Client) Publish(msgType string, data any) error {
	c.mutex.Lock()
	ch := c.channel
	c.mutex.Unlock()

	if ch == nil {
		return errors.New("not connected")
	}

	bodyData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := Message{
		Type: msgType,
		Data: bodyData,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",
		c.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (c *Client) Consume(
	ctx context.Context,
	handler func(msgType string, data []byte) error,
) error {
	c.mutex.Lock()
	ch := c.channel
	c.mutex.Unlock()

	if ch == nil {
		return errors.New("not connected")
	}

	deliveries, err := c.channel.Consume(
		c.queue,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				c.logger.Error(ctx.Err(), "rabbitmq - Consume - ctx.Done")
				return

			case err := <-c.closeChan:
				c.logger.Error(fmt.Errorf("connection closed: %w", err), "rabbitmq - Consume")
				return

			case d, ok := <-deliveries:
				if !ok {
					c.logger.Error(errors.New("delivery channel closed"), "rabbitmq - Consume")
					return
				}

				var msg Message
				if err := json.Unmarshal(d.Body, &msg); err != nil {
					d.Nack(false, false)
					continue
				}

				if err := handler(msg.Type, msg.Data); err != nil {
					d.Nack(false, true) // requeue
					continue
				}

				d.Ack(false)
			}
		}
	}()
	return nil
}

func (c *Client) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
