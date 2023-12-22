package server

import (
	"encoding/json"
	"errors"
	"fmt"
)

type EventJson struct {
	EventType string      `json:"type"`
	Payload   JsonPayload `json:"payload"`
}

func (e *EventJson) ParseMessage() (Message, error) {
	jsonMessage, _ := e.GetPayload().Serialize()
	switch e.GetType() {
	case EventMessageIn:
		var message *MessageInJson
		if err := json.Unmarshal(jsonMessage, &message); err != nil {
			return nil, errors.Join(ErrBadJsonPayload, fmt.Errorf("%v", err))
		}
		return message, nil
	case EventRoomCreate:
		var message *MessageCreateRoomJson
		if err := json.Unmarshal(jsonMessage, &message); err != nil {
			return nil, errors.Join(ErrBadJsonPayload, fmt.Errorf("%v", err))
		}
		return message, nil
	case EventRoomJoin:
		var message *MessageJoinRoomJson
		if err := json.Unmarshal(jsonMessage, &message); err != nil {
			return nil, errors.Join(ErrBadJsonPayload, fmt.Errorf("%v", err))
		}
		return message, nil
	case EventRoomGet:
		var message *MessageGetRoomJson
		if err := json.Unmarshal(jsonMessage, &message); err != nil {
			return nil, errors.Join(ErrBadJsonPayload, fmt.Errorf("%v", err))
		}
		return message, nil

	default:
		return nil, ErrUnsupportedMessageType
	}
}

func parseJsonEvent(payload []byte) (Event, error) {
	var event *EventJson
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, err
	}
	return event, nil
}
