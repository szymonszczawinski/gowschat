package server

import (
	"errors"
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

func createMessageOut(m Message, peerType PeerType) (Message, error) {
	switch v := m.(type) {
	case MessageIn:
		switch peerType {
		case PeerTypeJson:
			return v.GenerateMessageOut(), nil
		default:
			return nil, ErrUnsupporterPeerType
		}
	default:
		return m, nil

	}
}

func createMessageRoomList(p *ChatPeer) (Message, error) {
	rooms := []string{}
	for room := range p.server.rooms {
		rooms = append(rooms, room.name)
	}
	switch p.peerType {
	case PeerTypeJson:

		message := &MessageRoomListJson{
			Rooms: rooms,
		}
		return message, nil
	default:
		return nil, ErrUnsupporterPeerType
	}
}

func createMessageRoom(r *ChatRoom, peerType PeerType) (Message, error) {
	peers := []string{}
	for peer := range r.peers {
		peers = append(peers, peer.peerId)
	}
	switch peerType {
	case PeerTypeJson:
		message := &MessageRoomJson{
			RoomName: r.name,
			Peers:    peers,
		}
		return message, nil
	default:
		return nil, ErrUnsupporterPeerType
	}
}
