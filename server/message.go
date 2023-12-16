package server

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

var ErrUnsupportedMessageType = errors.New("unsupported message type")

type Event interface {
	GetType() string
	GetPayload() any
}

type EventJson struct {
	MessageType string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
}

func (e EventJson) GetType() string {
	return e.MessageType
}

func (e EventJson) GetPayload() any {
	return e.Payload
}

type EventHandler func(event Event, p *ChatPeer) error

const (
	EventSendMessage = "send_message"
)

type SentMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

func parseEvent(messageType int, payload []byte) (Event, error) {
	switch messageType {
	case websocket.TextMessage:

		var event EventJson
		if err := json.Unmarshal(payload, &event); err != nil {
			log.Printf("error marshalling message: %v", err)
			return nil, err
		}
		return event, nil
	case websocket.BinaryMessage:
		return nil, ErrEventNotSupported
	default:
		return nil, ErrUnsupportedMessageType
	}
}
