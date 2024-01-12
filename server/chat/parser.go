package chat

import (
	"encoding/json"
	"errors"
	"gowschat/server/api"
	"gowschat/server/chat/messages"

	"github.com/gorilla/websocket"
)

var (
	ErrorUnsupportedMessageType = errors.New("unsupported message type")
	ErrorJsonParseNotAMap       = errors.New("JSON element is not a map")
	ErrorJsonParseError         = errors.New("JSON parse error")
)

const (
	Headers         = "HEADERS"
	HeaderHxRequest = "HX-Request"
	True            = "true"
	MessageKey      = "chat_message"
	APIClent        = "apiClient"
	EventType       = "eventType"
	JSONPayloadData = "data"
	JSONMessageIM   = "message"
)

func ParseIncomingMessage(messageType int, payload []byte) (messages.Event, error) {
	switch messageType {
	case websocket.BinaryMessage:
		return messages.Event{}, ErrorUnsupportedMessageType

	case websocket.TextMessage:
		return parseTextMessage(payload)
	default:
		return messages.Event{}, ErrorUnsupportedMessageType
	}
}

func parseTextMessage(payload []byte) (messages.Event, error) {
	var payloadAsMap map[string]any
	json.Unmarshal(payload, &payloadAsMap)
	if header, ok := payloadAsMap[Headers]; ok {
		if headerMap, isMap := header.(map[string]any); isMap {
			if _, isHtmx := headerMap[HeaderHxRequest]; isHtmx {
				messageContent := payloadAsMap[MessageKey]
				if stringContent, isString := messageContent.(string); isString {
					return createEventFromText(stringContent)
				}
			}
		}
	} else if _, isAPIClinet := payloadAsMap[APIClent]; isAPIClinet {
		return createEventFromJson(payloadAsMap)
	}
	return messages.Event{}, ErrorUnsupportedMessageType
}

func createEventFromText(m string) (messages.Event, error) {
	return messages.NewEvent(api.EventMessageIM, messages.NewMessageIM(m)), nil
}

func createEventFromJson(payloadAsMap map[string]any) (messages.Event, error) {
	eventType, _ := payloadAsMap[EventType].(string)
	switch api.EventType(eventType) {
	case api.EventMessageIM:
		return createEventFromJsonMessageIM(payloadAsMap)
	}
	return messages.Event{}, ErrorJsonParseError
}

func createEventFromJsonMessageIM(payloadAsMap map[string]any) (messages.Event, error) {
	data := payloadAsMap[JSONPayloadData]
	if dataMap, isMap := data.(map[string]any); isMap {
		message := dataMap[JSONMessageIM]
		if stringMessage, isString := message.(string); isString {
			return messages.NewEvent(api.EventMessageIM, messages.NewMessageIM(stringMessage)), nil
		}
	}
	return messages.Event{}, ErrorJsonParseError
}

func toMap(a any) (map[string]any, error) {
	if dataMap, isMap := a.(map[string]any); isMap {
		return dataMap, nil
	}
	return nil, ErrorJsonParseNotAMap
}
