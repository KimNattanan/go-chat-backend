package rabbitmq

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func NewMessage(msgType string, data any) ([]byte, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	msg := Message{
		Type: msgType,
		Data: payload,
	}

	return json.Marshal(msg)
}

func ParseMessage(body []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(body, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

type Publisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

func NewPublisher(url, exchange string) (*Publisher, error) {
	if exchange == "" {
		exchange = _defaultExchangeName
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	if err := ch.ExchangeDeclare(
		exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}

	return &Publisher{
		conn:     conn,
		channel:  ch,
		exchange: exchange,
	}, nil
}

func (p *Publisher) Close() error {
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

func (p *Publisher) Publish(msgType string, data any) error {
	body, err := NewMessage(msgType, data)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		p.exchange,
		"",    // routing key ignored for fanout
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
