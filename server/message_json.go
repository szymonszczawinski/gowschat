package server

import (
	"encoding/json"
	"errors"
	"gowschat/server/api"
	"time"
)

type (
	JsonPayload []byte

	MessageInJson struct {
		Message string `json:"message"`
		From    string `json:"from"`
	}

	MessageOutJson struct {
		Sent time.Time `json:"sent"`
		MessageInJson
	}

	MessageErrorJson struct {
		Error string `json:"error"`
	}

	MessageRoomJson struct {
		RoomName string   `json:"roomName"`
		Peers    []string `json:"peers"`
	}

	MessageGetRoomJson struct {
		RoomName string `json:"roomName"`
	}
	MessageJoinRoomJson struct {
		RoomName string `json:"roomName"`
	}
	MessageRoomListJson struct {
		Rooms []string `json:"rooms"`
	}
	MessageCreateRoomJson struct {
		RoomName string `json:"roomName"`
	}
	MessageUserRegisterJson struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)

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
	return api.EventUnknown
}

func (m JsonPayload) CreateEvent(payload []byte) (api.Event, error) {
	return nil, api.ErrOperationNotSupported
}

func (e EventJson) GetType() string {
	return e.EventType
}

func (e EventJson) GetPayload() api.MessageSerializable {
	return e.Payload
}

func (e EventJson) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (m MessageInJson) GetMessage() string {
	return m.Message
}

func (m MessageInJson) GetFrom() string {
	return m.From
}

func (m MessageInJson) GetType() string {
	return api.EventMessageIn
}

func (m MessageOutJson) GetType() string {
	return api.EventMessageOut
}

func (m *MessageOutJson) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MessageOutJson) CreateEvent(payload []byte) (api.Event, error) {
	return &EventJson{
		EventType: m.GetType(),
		Payload:   payload,
	}, nil
}

func (m MessageRoomListJson) GetType() string {
	return api.EventRoomList
}

func (m *MessageRoomListJson) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MessageRoomListJson) CreateEvent(payload []byte) (api.Event, error) {
	return &EventJson{
		EventType: m.GetType(),
		Payload:   payload,
	}, nil
}

func (m MessageCreateRoomJson) GetType() string {
	return api.EventRoomCreate
}

func (m MessageCreateRoomJson) GetRoomName() string {
	return m.RoomName
}

func (m MessageJoinRoomJson) GetType() string {
	return api.EventRoomJoin
}

func (m MessageJoinRoomJson) GetRoomName() string {
	return m.RoomName
}

func (m MessageGetRoomJson) GetType() string {
	return api.EventRoomGet
}

func (m MessageGetRoomJson) GetRoomName() string {
	return m.RoomName
}

func (m MessageRoomJson) GetType() string {
	return api.EventRoom
}

func (m *MessageRoomJson) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MessageRoomJson) CreateEvent(payload []byte) (api.Event, error) {
	return &EventJson{
		EventType: m.GetType(),
		Payload:   payload,
	}, nil
}

func (m MessageErrorJson) GetType() string {
	return api.EventError
}

func (m *MessageErrorJson) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MessageErrorJson) CreateEvent(payload []byte) (api.Event, error) {
	return &EventJson{
		EventType: m.GetType(),
		Payload:   payload,
	}, nil
}

func (m *MessageUserRegisterJson) GetType() string {
	return api.EventChatRegister
}

func (m *MessageUserRegisterJson) GetEmail() string {
	return m.Email
}

func (m *MessageUserRegisterJson) GetPassword() string {
	return m.Password
}
