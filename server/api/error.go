package api

import "errors"

var (
	ErrOperationNotSupported = NewAppError("operation not suuported")
	// payload errors

	ErrIncorrectJsonFormat = NewAppError("incorrect JSON format")
	ErrBadJsonPayload      = NewAppError("bad payload in request")

	// messate errors

	ErrUnsupportedMessageType = NewAppError("unsupported message type")
	ErrIncorrectMessageTypeIn = NewAppError("message type incorrect for create output type")

	// event errors

	ErrEventNotSupported = NewAppError("this event type is not supported")

	// peer errors

	ErrUnsupporterPeerType = NewAppError("unsupported peer type")
	// room errors

	ErrRoomExist              = NewChatError("room witch such name exist")
	ErrRoomNotExist           = NewChatError("room does not exist")
	ErrRoomAlreadyJoined      = NewChatError("room already joined")
	ErrUserAlreadyRegisterred = NewChatError("user already registerred")
	ErrUserNotRegisterred     = NewChatError("user not registerred")

	ErrorJsonParseNotAMap    = errors.New("JSON element is not a map")
	ErrorJsonParseError      = errors.New("JSON parse error")
	ErrorJsonSerlializeError = errors.New("JSON serialize error")
)

type (
	// internal application AppError
	// not sent to the peers
	AppError struct {
		message string
	}

	// error being result of a stadard behavior, e.g. validation, etc
	// send back to peer as a Event/Message
	ChatError struct {
		message string
	}
)

func NewAppError(message string) error {
	return AppError{
		message: message,
	}
}

func NewChatError(message string) error {
	return ChatError{
		message: message,
	}
}

func (e AppError) Error() string {
	return e.message
}

func (e ChatError) Error() string {
	return e.message
}
