package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type JsonMessage []byte

func (m JsonMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *JsonMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("JsonMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

func (m JsonMessage) Serialize() ([]byte, error) {
	return []byte(m), nil
}

type EventJson struct {
	EventType string      `json:"type"`
	Payload   JsonMessage `json:"payload"`
}

func (e EventJson) GetType() string {
	return e.EventType
}

func (e EventJson) GetPayload() Message {
	return e.Payload
}

func (e EventJson) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

type SendMessageJson struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

func (m SendMessageJson) GetMessage() string {
	return m.Message
}

func (m SendMessageJson) GetFrom() string {
	return m.From
}

func (m SendMessageJson) GenerateOutMessage() OutMessage {
	return OutMessageJson{
		Sent: time.Now(),
		SendMessageJson: SendMessageJson{
			Message: m.GetMessage(),
			From:    m.GetFrom(),
		},
	}
}

type OutMessageJson struct {
	Sent time.Time `json:"sent"`
	SendMessageJson
}

func (m OutMessageJson) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func (m OutMessageJson) Wrap() Event {
	data, _ := m.Serialize()
	return EventJson{
		EventType: EventNewMessage,
		Payload:   data,
	}
}

type InitialMessageJson struct {
	PeerId string `json:"peerid"`
}

func (m InitialMessageJson) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func parseJsonMessage(e EventJson) (SendMessage, error) {
	switch e.GetType() {
	case EventSendMessage:
		var message SendMessageJson
		jsonMessage, _ := e.GetPayload().Serialize()
		if err := json.Unmarshal(jsonMessage, &message); err != nil {
			return nil, errors.Join(ErrBadJsonPayload, fmt.Errorf("%v", err))
		}
		return message, nil

	default:
		return nil, ErrUnsupportedMessageType
	}
}
