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
	MessageIn interface {
		GenerateMessageOut() Message
	}
	MessageCreateRoom interface {
		GetRoomName() string
	}
	MessageJoinRoom interface {
		GetRoomName() string
	}
	MessageLeaveRoom interface {
		GetRoomName() string
	}
	MessageGetRoom interface {
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

func createMessageOut(m Message) (Message, error) {
	if inMessage, ok := m.(MessageIn); ok {
		return inMessage.GenerateMessageOut(), nil
	}
	return nil, ErrUnsupportedMessageType
}

func createMessageRoomList(p *ChatPeer) Message {
	rooms := []string{}
	for room := range p.server.rooms {
		rooms = append(rooms, room.name)
	}
	message := &MessageRoomListJson{
		Rooms: rooms,
	}
	return message
}

func createMessageRoom(r *ChatRoom) Message {
	peers := []string{}
	for peer := range r.peers {
		peers = append(peers, peer.peerId)
	}
	message := &MessageRoomJson{
		RoomName: r.name,
		Peers:    peers,
	}
	return message
}
