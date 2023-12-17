package server

import (
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
	EventCreateRoomAck      = "create_room_ack"
	EventGetRoom            = "get_room"
	EventJoinRoom           = "join_room"
	EventJoinRoomAck        = "join_room_ack"
	EventLeaveRoom          = "leave_room"
	EventLeaveRoomAck       = "leave_room_ack"
	EventGetRoomList        = "list_rooms"
	EventRoomList           = "list_rooms_resp"
	EventError              = "error"
	EventUnknown            = "unknown"
)

type (
	Event interface {
		GetType() string
		GetPayload() Message
		Serialize() ([]byte, error)
		ParseMessage() (Message, error)
	}
)
type EventHandler func(event Event, p *ChatPeer) error

func parseEvent(messageType int, payload []byte) (Event, error) {
	switch messageType {
	case websocket.TextMessage:
		return parseJsonEvent(payload)
	case websocket.BinaryMessage:
		return nil, ErrUnsupportedMessageType
	default:
		return nil, ErrEventNotSupported
	}
}

func createEvent(m Message, peerType PeerType) (Event, error) {
	messagedata, err := m.Serialize()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if peerType == PeerTypeJson {
		return &EventJson{
			EventType: m.GetType(),
			Payload:   messagedata,
		}, nil
	}
	if peerType == PeerTypeProto {
		return nil, errors.Join(ErrUnsupporterPeerType, fmt.Errorf("%v", peerType))
	}
	return nil, ErrUnsupporterPeerType
}

func HandlerMessageIn(event Event, p *ChatPeer) error {
	if message, err := event.ParseMessage(); err != nil {
		log.Println("ERROR", err)
	} else {
		// Broadcast to all other Clients
		for peer := range p.server.peers {
			// FIXME:: re-enable same peer check
			// if p.peerId != peer.peerId {
			if peer.status == PeerStatusOnline {
				peer.outgoing <- message
			}
		}
	}
	return nil
}

func HandlerCreateRoom(e Event, p *ChatPeer) error {
	if message, err := e.ParseMessage(); err != nil {
		log.Println("ERROR", err)
	} else {
		if createRoomMessage, ok := message.(MessageCreateRoom); ok {
			if err := p.server.createRoom(createRoomMessage.GetRoomName(), p); err != nil {
				return err
			} else {
				log.Println("room created", createRoomMessage.GetRoomName())
			}
		} else {
			return ErrUnsupportedMessageType
		}
	}
	return nil
}

func HandlerJoinRoom(e Event, p *ChatPeer) error {
	if message, err := e.ParseMessage(); err != nil {
		log.Println("ERROR", err)
	} else {
		if joinRoomMessage, ok := message.(MessageJoinRoom); ok {
			room, err := p.server.joinRoom(joinRoomMessage.GetRoomName(), p)
			if err != nil {
				return err
			}
			message, err := createMessageRoom(room, p.peerType)
			if err != nil {
				log.Println(err)
				return err
			}
			p.outgoing <- message

		} else {
			return ErrUnsupportedMessageType
		}
	}

	return nil
}

func HandlerGetRoom(e Event, p *ChatPeer) error {
	if message, err := e.ParseMessage(); err != nil {
		log.Println("ERROR", err)
	} else {
		if getRoomMessage, ok := message.(MessageGetRoom); ok {
			room, err := p.server.getRoom(getRoomMessage.GetRoomName())
			if err != nil {
				return err
			}
			message, err := createMessageRoom(room, p.peerType)
			if err != nil {
				log.Println(err)
				return err
			}
			p.outgoing <- message

		} else {
			return ErrUnsupportedMessageType
		}
	}

	return nil
}

func HandlerLeaveRoom(e Event, p *ChatPeer) error {
	// TODO: todo
	return nil
}

func HandlerRoomList(e Event, p *ChatPeer) error {
	message, err := createMessageRoomList(p)
	if err != nil {
		log.Println(err)
		return err
	}
	p.outgoing <- message
	return nil
}
