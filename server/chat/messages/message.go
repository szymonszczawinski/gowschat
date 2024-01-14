package messages

import "gowschat/server/api"

type (
	Event struct {
		message   api.IMessage
		eventType api.EventType
	}
	MessageIM struct {
		message string
		from    string
	}
	MessageError struct {
		errorMessage string
	}
)

func NewEvent(m api.IMessage) Event {
	return Event{
		eventType: api.MessageTypeToEventType[m.GetType()],
		message:   m,
	}
}

func NewEventWithType(etype api.EventType, m api.IMessage) Event {
	return Event{
		eventType: etype,
		message:   m,
	}
}

func (e Event) GetType() api.EventType {
	return e.eventType
}

func (e Event) GetMessage() api.IMessage {
	return e.message
}

func NewMessageIM(m string, from string) api.IMessage {
	return MessageIM{
		message: m,
		from:    from,
	}
}

func (m MessageIM) GetType() api.MesageType {
	return api.MessageIM
}

func (m MessageIM) GetFrom() string {
	return m.from
}

func (m MessageIM) GetMessage() string {
	return m.message
}

func NewMessageError(e error) api.IMessage {
	return MessageError{
		errorMessage: e.Error(),
	}
}

func (m MessageError) GetType() api.MesageType {
	return api.MessageError
}

func (m MessageError) GetErrorMessage() string {
	return m.errorMessage
}
