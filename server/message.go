package server

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrUnsupportedMessageType = errors.New("unsupported message type")
	ErrIncorrectJsonFormat    = errors.New("incorrect JSON format")
	ErrBadJsonPayload         = errors.New("bad payload in request")
)

type (
	Message interface {
		GetType() string
		Serialize() ([]byte, error)
	}
	InMessage interface {
		GenerateOutMessage() Message
	}
	CreateRoomMessage interface {
		GetRoomName() string
	}
)

func parseMessage(e Event, peerType PeerType) (Message, error) {
	switch peerType {
	case PeerTypeJson:
		return parseJsonMessage(e.GetType(), e.GetPayload())
	default:
		return nil, ErrUnsupportedMessageType
	}
}

func serializeMessage(out Message, peerType PeerType) ([]byte, error) {
	if peerType == PeerTypeJson {
		return json.Marshal(out)
	}
	if peerType == PeerTypeProto {
		return nil, errors.Join(ErrUnsupporterPeerType, fmt.Errorf("%v", peerType))
	}
	return nil, ErrUnsupporterPeerType
}

func createOutMessage(m Message) (Message, error) {
	if inMessage, ok := m.(InMessage); ok {
		return inMessage.GenerateOutMessage(), nil
	}
	return nil, ErrUnsupportedMessageType
}

func createRoomListMessage(p *ChatPeer) Message {
	rooms := []string{}
	for room := range p.server.rooms {
		rooms = append(rooms, room.name)
	}
	message := &ListRoomsMessageRespJson{
		Rooms: rooms,
	}
	return message
}
