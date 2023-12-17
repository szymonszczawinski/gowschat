package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type JsonPayload []byte

func (m JsonPayload) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *JsonPayload) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("JsonMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

func (m JsonPayload) Serialize() ([]byte, error) {
	return []byte(m), nil
}

func (m JsonPayload) GetType() string {
	return EventUnknown
}

type EventJson struct {
	EventType string      `json:"type"`
	Payload   JsonPayload `json:"payload"`
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

type MessageBaseJson struct{}

func (m *MessageBaseJson) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

type MessageInJson struct {
	MessageBaseJson
	Message string `json:"message"`
	From    string `json:"from"`
}

func (m MessageInJson) GetMessage() string {
	return m.Message
}

func (m MessageInJson) GetFrom() string {
	return m.From
}

func (m MessageInJson) GetType() string {
	return EventInMessage
}

func (m MessageInJson) GenerateOutMessage() Message {
	return &MessageOutJson{
		Sent: time.Now(),
		MessageInJson: MessageInJson{
			Message: m.GetMessage(),
			From:    m.GetFrom(),
		},
	}
}

type MessageOutJson struct {
	Sent time.Time `json:"sent"`
	MessageInJson
}

func (m MessageOutJson) GetType() string {
	return EventOutMessage
}

type MessageRoomListJson struct {
	MessageBaseJson
	Rooms []string `json:"rooms"`
}

func (m MessageRoomListJson) GetType() string {
	return EventRoomList
}

type MessageCreateRoomJson struct {
	MessageBaseJson
	RoomName string `json:"roomName"`
}

func (m MessageCreateRoomJson) GetType() string {
	return EventCreateRoom
}

func (m MessageCreateRoomJson) GetRoomName() string {
	return m.RoomName
}

type MessageJoinRoomJson struct {
	MessageBaseJson
	RoomName string `json:"roomName"`
}

func (m MessageJoinRoomJson) GetType() string {
	return EventJoinRoom
}

func (m MessageJoinRoomJson) GetRoomName() string {
	return m.RoomName
}

type MessageGetRoomJson struct {
	MessageBaseJson
	RoomName string `json:"roomName"`
}

func (m MessageGetRoomJson) GetType() string {
	return EventGetRoom
}

func (m MessageGetRoomJson) GetRoomName() string {
	return m.RoomName
}

type MessageRoomJson struct {
	MessageBaseJson
	RoomName string   `json:"roomName"`
	Peers    []string `json:"peers"`
}

func (m MessageRoomJson) GetType() string {
	return EventJoinRoomAck
}

func parseJsonEvent(payload []byte) (Event, error) {
	var event *EventJson
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Printf("error marshalling message: %v", err)
		return nil, err
	}
	return event, nil
}

func (e *EventJson) ParseMessage() (Message, error) {
	jsonMessage, _ := e.GetPayload().Serialize()
	switch e.GetType() {
	case EventInMessage:
		var message *MessageInJson
		if err := json.Unmarshal(jsonMessage, &message); err != nil {
			return nil, errors.Join(ErrBadJsonPayload, fmt.Errorf("%v", err))
		}
		return message, nil
	case EventCreateRoom:
		var message *MessageCreateRoomJson
		if err := json.Unmarshal(jsonMessage, &message); err != nil {
			return nil, errors.Join(ErrBadJsonPayload, fmt.Errorf("%v", err))
		}
		return message, nil
	case EventJoinRoom:
		var message *MessageJoinRoomJson
		if err := json.Unmarshal(jsonMessage, &message); err != nil {
			return nil, errors.Join(ErrBadJsonPayload, fmt.Errorf("%v", err))
		}
		return message, nil
	case EventGetRoom:
		var message *MessageGetRoomJson
		if err := json.Unmarshal(jsonMessage, &message); err != nil {
			return nil, errors.Join(ErrBadJsonPayload, fmt.Errorf("%v", err))
		}
		return message, nil

	default:
		return nil, ErrUnsupportedMessageType
	}
}
