package chatjson

//
// import (
// 	"encoding/json"
// 	"gowschat/server/api"
//
// 	"github.com/gorilla/websocket"
// )
//
// type ParserJson struct{}
//
// func (p ParserJson) ParseEvent(messageType int, payload []byte) (api.Event, error) {
// 	if messageType == websocket.TextMessage {
// 		var event *EventJson
// 		if err := json.Unmarshal(payload, &event); err != nil {
// 			return nil, err
// 		}
// 		return event, nil
// 	}
// 	return nil, api.ErrUnsupportedMessageType
// }
