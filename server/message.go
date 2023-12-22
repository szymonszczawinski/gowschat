package server

type (
	Message interface {
		GetType() string
	}
	MessageSerializable interface {
		Message
		Serialize() ([]byte, error)
	}
	MessageIn interface {
		GenerateMessageOut() MessageSerializable
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
	MessageRegister interface {
		GetEmail() string
		GetPassword() string
	}
)

func createMessageOut(m Message, peerType PeerType) (MessageSerializable, error) {
	switch v := m.(type) {
	case MessageIn:
		switch peerType {
		case PeerTypeJson:
			return v.GenerateMessageOut(), nil
		default:
			return nil, ErrUnsupporterPeerType
		}
	default:
		return nil, ErrIncorrectMessageTypeIn

	}
}

func createMessageRoomList(p *ChatPeer) (MessageSerializable, error) {
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

func createMessageRoom(r *ChatRoom, peerType PeerType) (MessageSerializable, error) {
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

func createMessageError(appError error, peerType PeerType) (MessageSerializable, error) {
	switch peerType {
	case PeerTypeJson:
		message := &MessageErrorJson{
			Error: appError.Error(),
		}
		return message, nil
	default:
		return nil, ErrUnsupporterPeerType
	}
}
