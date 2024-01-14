package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"gowschat/server/api"
	"gowschat/server/chat/messages"
	"log"

	"github.com/gorilla/websocket"
)

func ParseIncomingMessage(peer api.IChatPeer, messageType int, payload []byte) (messages.Event, error) {
	switch messageType {
	case websocket.BinaryMessage:
		return messages.Event{}, api.ErrUnsupportedMessageType
	case websocket.TextMessage:
		return parseTextMessage(peer, payload)
	default:
		return messages.Event{}, api.ErrUnsupportedMessageType
	}
}

func parseTextMessage(peer api.IChatPeer, payload []byte) (messages.Event, error) {
	log.Println("parseTextMessage:", string(payload))
	var payloadAsMap map[string]any
	json.Unmarshal(payload, &payloadAsMap)
	switch peer.GetType() {
	case api.PeerTypeJson:
		return createEventFromJson(peer, payloadAsMap)
	case api.PeerTypeWeb:
		return createEventFromText(peer, payloadAsMap)
	default:
		return messages.Event{}, api.ErrUnsupporterPeerType
	}
}

func isHtmxPeer(payloadAsMap map[string]any) bool {
	if header, ok := payloadAsMap[api.JSONHeaders]; ok {
		if headerMap, isMap := header.(map[string]any); isMap {
			if _, isHtmx := headerMap[api.JSONHeadersHxRequest]; isHtmx {
				return true
			}
		}
	}
	return false
}

func isApiPeer(payloadAsMap map[string]any) bool {
	_, isAPIClinet := payloadAsMap[api.JSONAPIClient]

	return isAPIClinet
}

func createEventFromText(peer api.IChatPeer, payloadAsMap map[string]any) (messages.Event, error) {
	messageContent := payloadAsMap[api.JSONWebMessage]
	if stringContent, isString := messageContent.(string); isString {
		return messages.NewEvent(api.EventMessageIM, messages.NewMessageIM(stringContent, peer.GetUserName())), nil
	}
	return messages.Event{}, errors.Join(api.ErrorJsonParseError, fmt.Errorf("no a string message content %v -> %T", messageContent, messageContent))
}

func createEventFromJson(peer api.IChatPeer, payloadAsMap map[string]any) (messages.Event, error) {
	eventType, _ := payloadAsMap[api.JSONEventType].(string)
	switch api.EventType(eventType) {
	case api.EventMessageIM:
		return createEventFromJsonMessageIM(peer, payloadAsMap)
	}
	return messages.Event{}, api.ErrorJsonParseError
}

func createEventFromJsonMessageIM(peer api.IChatPeer, payloadAsMap map[string]any) (messages.Event, error) {
	data := payloadAsMap[api.JSONData]
	if dataMap, isMap := data.(map[string]any); isMap {
		message := dataMap[api.JSONDataMessage]
		if stringMessage, isString := message.(string); isString {
			return messages.NewEvent(api.EventMessageIM, messages.NewMessageIM(stringMessage, peer.GetUserName())), nil
		}
	}
	return messages.Event{}, api.ErrorJsonParseError
}

func toMap(a any) (map[string]any, error) {
	if dataMap, isMap := a.(map[string]any); isMap {
		return dataMap, nil
	}
	return nil, api.ErrorJsonParseNotAMap
}
