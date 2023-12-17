package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

const (
	// EventInMessage is the event name for new chat messages sent
	EventInMessage = "in_message"
	// EventOutMessage is a response to send_message
	EventOutMessage = "out_message"
	// EventInitialPeerMessage is a initial message sent to peer after connect
	EventInitialPeerMessage = "initial_message"
	EventCreateRoom         = "create_room"
	EventJoinRoom           = "join_room"
	EventLeaveRoom          = "leave_room"
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
	InMessage interface {
		GetMessage() string
		GetFrom() string
		GenerateOutMessage() OutMessage
	}

	OutMessage interface {
		Serialize() ([]byte, error)
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

func parseMessage(e Event, peerType PeerType) (InMessage, error) {
	switch peerType {
	case PeerTypeJson:
		return parseJsonMessage(e.GetType(), e.GetPayload())
	default:
		return nil, ErrUnsupportedMessageType
	}
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

func createOutEvent(out OutMessage, peerType PeerType) (Event, error) {
	messagedata, err := serializeOutMessage(out, peerType)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if peerType == PeerTypeJson {
		return EventJson{
			EventType: EventOutMessage,
			Payload:   messagedata,
		}, nil
	}
	if peerType == PeerTypeProto {
		return nil, errors.Join(ErrUnsupporterPeerType, fmt.Errorf("%v", peerType))
	}
	return nil, ErrUnsupporterPeerType
}
