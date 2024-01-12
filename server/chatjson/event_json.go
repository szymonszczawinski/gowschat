package chatjson

//
// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"gowschat/server/api"
// )
//
// type EventJson struct {
// 	EventType string      `json:"type"`
// 	Payload   JsonPayload `json:"payload"`
// }
//
// func (e *EventJson) ParseMessage() (api.Message, error) {
// 	jsonMessage, _ := e.GetPayload().Serialize()
// 	switch e.GetType() {
// 	case api.EventMessageIn:
// 		var message *MessageInJson
// 		if err := json.Unmarshal(jsonMessage, &message); err != nil {
// 			return nil, errors.Join(api.ErrBadJsonPayload, fmt.Errorf("%v", err))
// 		}
// 		return message, nil
// 	case api.EventRoomCreate:
// 		var message *MessageCreateRoomJson
// 		if err := json.Unmarshal(jsonMessage, &message); err != nil {
// 			return nil, errors.Join(api.ErrBadJsonPayload, fmt.Errorf("%v", err))
// 		}
// 		return message, nil
// 	case api.EventRoomJoin:
// 		var message *MessageJoinRoomJson
// 		if err := json.Unmarshal(jsonMessage, &message); err != nil {
// 			return nil, errors.Join(api.ErrBadJsonPayload, fmt.Errorf("%v", err))
// 		}
// 		return message, nil
// 	case api.EventRoomGet:
// 		var message *MessageGetRoomJson
// 		if err := json.Unmarshal(jsonMessage, &message); err != nil {
// 			return nil, errors.Join(api.ErrBadJsonPayload, fmt.Errorf("%v", err))
// 		}
// 		return message, nil
//
// 	default:
// 		return nil, api.ErrUnsupportedMessageType
// 	}
// }
//
// func parseJsonEvent(payload []byte) (api.Event, error) {
// 	var event *EventJson
// 	if err := json.Unmarshal(payload, &event); err != nil {
// 		return nil, err
// 	}
// 	return event, nil
// }
