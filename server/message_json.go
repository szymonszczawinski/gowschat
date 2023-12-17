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

func (m JsonMessage) GetType() string {
	return EventUnknown
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

type BaseMessageJson struct{}

func (m *BaseMessageJson) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

type InMessageJson struct {
	BaseMessageJson
	Message string `json:"message"`
	From    string `json:"from"`
}

func (m InMessageJson) GetMessage() string {
	return m.Message
}

func (m InMessageJson) GetFrom() string {
	return m.From
}

func (m InMessageJson) GetType() string {
	return EventInMessage
}

func (m InMessageJson) GenerateOutMessage() Message {
	return &OutMessageJson{
		Sent: time.Now(),
		InMessageJson: InMessageJson{
			Message: m.GetMessage(),
			From:    m.GetFrom(),
		},
	}
}

type OutMessageJson struct {
	Sent time.Time `json:"sent"`
	InMessageJson
}

func (m OutMessageJson) GetType() string {
	return EventOutMessage
}

type ListRoomsMessageRespJson struct {
	BaseMessageJson
	Rooms []string `json:"rooms"`
}

func (m ListRoomsMessageRespJson) GetType() string {
	return EventListRoomsResp
}

type CreateRoomMessageJson struct {
	BaseMessageJson
	RoomName string `json:"roomName"`
}

func (m CreateRoomMessageJson) GetType() string {
	return EventCreateRoom
}

func (m CreateRoomMessageJson) GetRoomName() string {
	return m.RoomName
}

func parseJsonMessage(eventType string, payload Message) (Message, error) {
	jsonMessage, _ := payload.Serialize()
	switch eventType {
	case EventInMessage:
		var message *InMessageJson
		if err := json.Unmarshal(jsonMessage, &message); err != nil {
			return nil, errors.Join(ErrBadJsonPayload, fmt.Errorf("%v", err))
		}
		return message, nil
	case EventCreateRoom:
		var message *CreateRoomMessageJson
		if err := json.Unmarshal(jsonMessage, &message); err != nil {
			return nil, errors.Join(ErrBadJsonPayload, fmt.Errorf("%v", err))
		}
		return message, nil

	default:
		return nil, ErrUnsupportedMessageType
	}
}
