package api

const (
	JSONHeaders          = "HEADERS"
	JSONHeadersHxRequest = "HX-Request"
	JSONValueTrue        = "true"
	JSONWebMessage       = "chat_message"
	JSONAPIClient        = "apiClient"
	JSONEventType        = "eventType"
	JSONData             = "data"
	JSONDataMessage      = "message"
	JSONDataFrom         = "from"
	JSONDataError        = "error"
)

type (
	IParser interface {
		ParseEvent(messageType int, pyload []byte) (IEvent, error)
	}

	ISerializer       interface{}
	IParserSerializer interface {
		IParser
		ISerializer
	}
)
