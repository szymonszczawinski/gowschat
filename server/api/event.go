package api

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
)
