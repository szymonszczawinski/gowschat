package messages

import "gowschat/server/api"

type (
	Event struct {
		message   api.IMessage
		eventType api.EventType
	}
	MessageIM struct {
		m string
	}
)

func NewEvent(etype api.EventType, m api.IMessage) Event {
	return Event{
		eventType: etype,
		message:   m,
	}
}

func NewMessageIM(m string) api.IMessage {
	return MessageIM{
		m: m,
	}
}

func (e Event) GetType() api.EventType {
	return e.eventType
}

func (e Event) GetMessage() api.IMessage {
	return e.message
}

func (m MessageIM) GetType() api.MesageType {
	return api.MessageIM
}
