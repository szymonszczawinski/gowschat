package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"gowschat/server/api"
	"gowschat/server/chat/messages"
)

func SerializeEvent(e messages.Event) ([]byte, error) {
	messageMap, err := serializeMessage(e.GetMessage())
	if err != nil {
		return nil, err
	}
	dataMap := map[string]any{}
	dataMap[api.JSONEventType] = e.GetType()
	dataMap[api.JSONData] = messageMap
	serilaized, err := json.Marshal(dataMap)
	return serilaized, err
}

func serializeMessage(m api.IMessage) (map[string]any, error) {
	switch m.GetType() {
	case api.MessageError:
		return seralizeMessageError(m)
	case api.MessageIM:
		return serializeMessageiIM(m)
	}
	return nil, api.ErrUnsupportedMessageType
}

func serializeMessageiIM(m api.IMessage) (map[string]any, error) {
	if imMessage, isIMMessage := m.(messages.MessageIM); isIMMessage {
		messageMap := map[string]any{}
		messageMap[api.JSONDataMessage] = imMessage.GetMessage()
		messageMap[api.JSONDataFrom] = imMessage.GetFrom()
		return messageMap, nil
	}
	return nil, errors.Join(api.ErrorJsonSerlializeError, fmt.Errorf("message has incorrect type %T", m))
}

func seralizeMessageError(m api.IMessage) (map[string]any, error) {
	if erroMessage, isErrorMessage := m.(messages.MessageError); isErrorMessage {
		messageMap := map[string]any{}
		messageMap[api.JSONDataError] = erroMessage.GetErrorMessage()
		return messageMap, nil
	}
	return nil, errors.Join(api.ErrorJsonSerlializeError, fmt.Errorf("message has incorrect type %T", m))
}
