package api

type EventType string

const (
	// EventMessageIM is the event name for new chat messages sent
	EventMessageIM    EventType = "message_im"
	EventChatRegister EventType = "chat_register"
	EventChatJoin     EventType = "chat_join"
	// EventMessageOut is a response to send_message
	EventMessageOut EventType = "message_out"
	// EventInitialPeerMessage is a initial message sent to peer after connect
	EventRoomCreate    EventType = "room_create"
	EventRoomCreateAck EventType = "room_create_ack"
	EventRoomGet       EventType = "room_get"
	EventRoomJoin      EventType = "room_join"
	EventRoom          EventType = "room"
	EventRoomLeave     EventType = "room_leave"
	EventRoomLeaveAck  EventType = "room_leave_ack"
	EventRoomListGet   EventType = "room_list_get"
	EventRoomList      EventType = "room_list"
	EventError         EventType = "error"
	EventUnknown       EventType = "unknown"
)

type MesageType string

const (
	MessageIM    MesageType = "IM"
	MessageError MesageType = "error"
)

type (
	IEvent interface {
		GetType() string
		GetPayload() IMessageSerializable
		Serialize() ([]byte, error)
		ParseMessage() (IMessage, error)
	}
	IMessage interface {
		GetType() MesageType
	}
	IMessageSerializable interface {
		IMessage
		Serialize() ([]byte, error)
		CreateEvent([]byte) (IEvent, error)
	}
	IMessageIM interface {
		GetFrom() string
		GetMessage() string
	}
	IMessageCreateRoom interface {
		GetRoomName() string
	}
	IMessageJoinRoom interface {
		GetRoomName() string
	}
	IMessageLeaveRoom interface {
		GetRoomName() string
	}
	IMessageGetRoom interface {
		GetRoomName() string
	}
	IMessageRegister interface {
		GetEmail() string
		GetPassword() string
	}
)
