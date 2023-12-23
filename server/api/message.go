package api

type (
	Message interface {
		GetType() string
	}
	MessageSerializable interface {
		Message
		Serialize() ([]byte, error)
		CreateEvent([]byte) (Event, error)
	}
	MessageIn interface {
		GetFrom() string
		GetMessage() string
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
