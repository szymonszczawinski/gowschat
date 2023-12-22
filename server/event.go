package server

import (
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

const (
	EventChatRegister = "chat_register"
	EventChatJoin     = "chat_join"
	// EventMessageIn is the event name for new chat messages sent
	EventMessageIn = "message_in"
	// EventMessageOut is a response to send_message
	EventMessageOut = "message_out"
	// EventInitialPeerMessage is a initial message sent to peer after connect
	EventRoomCreate    = "room_create"
	EventRoomCreateAck = "room_create_ack"
	EventRoomGet       = "room_get"
	EventRoomJoin      = "room_join"
	EventRoom          = "room"
	EventRoomLeave     = "room_leave"
	EventRoomLeaveAck  = "room_leave_ack"
	EventRoomListGet   = "room_list_get"
	EventRoomList      = "room_list"
	EventError         = "error"
	EventUnknown       = "unknown"
)

type (
	Event interface {
		GetType() string
		GetPayload() MessageSerializable
		Serialize() ([]byte, error)
		ParseMessage() (Message, error)
	}

	EventHandler func(event Event, p *ChatPeer) error
)

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

func createEvent(m MessageSerializable, peerType PeerType) (Event, error) {
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

func HandlerChatRegister(e Event, p *ChatPeer) error {
	if message, err := e.ParseMessage(); err != nil {
		log.Println("ERROR", err)
	} else {
		if messageRegister, ok := message.(MessageRegister); ok {
			if err := p.server.RegisterPeer(p, messageRegister.GetEmail(), messageRegister.GetPassword()); err != nil {
				return err
			}
		} else {
			return ErrUnsupportedMessageType
		}
	}
	return nil
}

func HandlerMessageIn(event Event, p *ChatPeer) error {
	if message, err := event.ParseMessage(); err != nil {
		log.Println("ERROR", err)
	} else {
		// Broadcast to all other Clients
		messageOut, err := createMessageOut(message, p.peerType)
		if err != nil {
			return err
		}
		for peer := range p.server.connectedPeers {
			// FIXME:: re-enable same peer check
			// if p.peerId != peer.peerId {
			if peer.status == PeerStatusOnline {
				peer.outgoing <- messageOut
			}
		}
	}
	return nil
}

func HandlerCreateRoom(e Event, p *ChatPeer) error {
	if message, parseErr := e.ParseMessage(); parseErr != nil {
		log.Println("ERROR", parseErr)
	} else {
		if createRoomMessage, ok := message.(MessageCreateRoom); ok {
			room, appErr := p.server.createRoom(createRoomMessage.GetRoomName(), p)
			if appErr != nil {
				return appErr
			}
			roomMessage, err := createMessageRoom(room, p.peerType)
			if err != nil {
				return err
			}
			p.outgoing <- roomMessage

		} else {
			return ErrUnsupportedMessageType
		}
	}
	return nil
}

func HandlerJoinRoom(e Event, p *ChatPeer) error {
	if message, parseErr := e.ParseMessage(); parseErr != nil {
		log.Println("ERROR", parseErr)
	} else {
		if joinRoomMessage, ok := message.(MessageJoinRoom); ok {
			room, appErr := p.server.joinRoom(joinRoomMessage.GetRoomName(), p)
			if appErr != nil {
				return appErr
			}
			message, err := createMessageRoom(room, p.peerType)
			if err != nil {
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
	if message, parseErr := e.ParseMessage(); parseErr != nil {
		log.Println("ERROR", parseErr)
	} else {
		if getRoomMessage, ok := message.(MessageGetRoom); ok {
			room, appErr := p.server.getRoom(getRoomMessage.GetRoomName())
			if appErr != nil {
				return appErr
			}
			message, appErr := createMessageRoom(room, p.peerType)
			if appErr != nil {
				return appErr
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
		return err
	}
	p.outgoing <- message
	return nil
}

func HandlerError(err error, e Event, p *ChatPeer) {
	log.Printf("ERROR :: handle event %v -> %v", e.GetType(), err.Error())
	var appErr ChatError
	if errors.As(err, &appErr) {
		errorMessage, e := createMessageError(err, p.peerType)
		if e != nil {
			log.Printf("ERROR :: handle application error %v -> %v", appErr, err)
		}

		p.outgoing <- errorMessage
	}
}
