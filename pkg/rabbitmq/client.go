package rabbitmq

import (
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Client struct {
	url   string
	queue string

	conn *amqp.Connection

	pubCh  *amqp.Channel
	consCh *amqp.Channel

	mu sync.RWMutex
}

const reconnectDelay = 5 * time.Second

func New(url, queue string) *Client {
	return &Client{
		url:   url,
		queue: queue,
	}
}

func (c *Client) Connect() error {
	for {
		conn, err := amqp.Dial(c.url)
		if err != nil {
			time.Sleep(reconnectDelay)
			continue
		}

		pubCh, err := conn.Channel()
		if err != nil {
			conn.Close()
			continue
		}

		consCh, err := conn.Channel()
		if err != nil {
			conn.Close()
			continue
		}

		_, err = consCh.QueueDeclare(
			c.queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			conn.Close()
			continue
		}

		c.mu.Lock()
		c.conn = conn
		c.pubCh = pubCh
		c.consCh = consCh
		c.mu.Unlock()

		return nil
	}
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pubCh != nil {
		c.pubCh.Close()
	}

	if c.consCh != nil {
		c.consCh.Close()
	}

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
