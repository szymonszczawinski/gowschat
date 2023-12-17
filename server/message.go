package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

const (
	// EventSendMessage is the event name for new chat messages sent
	EventSendMessage = "send_message"
	// EventNewMessage is a response to send_message
	EventNewMessage = "new_message"
	// EventInitialPeerMessage is a initial message sent to peer after connect
	EventInitialPeerMessage = "initial_message"
)

var (
	ErrUnsupportedMessageType = errors.New("unsupported message type")
	ErrIncorrectJsonFormat    = errors.New("incorrect JSON format")
	ErrBadJsonPayload         = errors.New("bad payload in request")
)

type (
	Event interface {
		GetType() string
		GetPayload() Message
		Serialize() ([]byte, error)
	}

	Message interface {
		Serialize() ([]byte, error)
	}
	SendMessage interface {
		GetMessage() string
		GetFrom() string
		GenerateOutMessage() OutMessage
	}

	OutMessage interface {
		Serialize() ([]byte, error)
		Wrap() Event
	}

	InitialMessage interface {
		SetPayload(p ChatPeer)
	}
)

type EventHandler func(event Event, p *ChatPeer) error

func parseEvent(messageType int, payload []byte) (Event, error) {
	switch messageType {
	case websocket.TextMessage:

		var event EventJson
		if err := json.Unmarshal(payload, &event); err != nil {
			log.Printf("error marshalling message: %v", err)
			return nil, err
		}
		return event, nil
	case websocket.BinaryMessage:
		return nil, ErrUnsupportedMessageType
	default:
		return nil, ErrEventNotSupported
	}
}

func parseSendMessage(e Event) (SendMessage, error) {
	switch v := e.(type) {
	case EventJson:
		return parseJsonMessage(EventJson(v))
	default:
		return nil, ErrUnsupportedMessageType
	}
}

func generateInitialEvent(p ChatPeer) Event {
	message := InitialMessageJson{PeerId: p.peerId}
	jsonMessage, _ := message.Serialize()
	event := EventJson{
		EventType: EventInitialPeerMessage,
		Payload:   jsonMessage,
	}
	return event
}

func serializeOutMessage(out OutMessage, peerType PeerType) ([]byte, error) {
	if peerType == PeerTypeJson {
		return json.Marshal(out)
	}
	if peerType == PeerTypeProto {
		return nil, errors.Join(ErrUnsupporterPeerType, fmt.Errorf("%v", peerType))
	}
	return nil, ErrUnsupporterPeerType
}

func wrapOutMessage(data []byte, peerType PeerType) (Event, error) {
	if peerType == PeerTypeJson {
		return EventJson{
			EventType: EventNewMessage,
			Payload:   data,
		}, nil
	}
	if peerType == PeerTypeProto {
		return nil, errors.Join(ErrUnsupporterPeerType, fmt.Errorf("%v", peerType))
	}
	return nil, ErrUnsupporterPeerType
}
