package rabbitmq

import "encoding/json"

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
