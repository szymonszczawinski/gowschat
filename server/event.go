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
	EventCreateRoomAck      = "create_room_ack"
	EventJoinRoom           = "join_room"
	EventJoinRoomAck        = "join_room_ack"
	EventLeaveRoom          = "leave_room"
	EventLeaveRoomAck       = "leave_room_ack"
	EventListRooms          = "list_rooms"
	EventListRoomsResp      = "list_rooms_resp"
	EventUnknown            = "unknown"
)

type (
	Event interface {
		GetType() string
		GetPayload() Message
		Serialize() ([]byte, error)
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

func createEvent(m Message, peerType PeerType) (Event, error) {
	messagedata, err := serializeMessage(m, peerType)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if peerType == PeerTypeJson {
		return EventJson{
			EventType: m.GetType(),
			Payload:   messagedata,
		}, nil
	}
	if peerType == PeerTypeProto {
		return nil, errors.Join(ErrUnsupporterPeerType, fmt.Errorf("%v", peerType))
	}
	return nil, ErrUnsupporterPeerType
}
