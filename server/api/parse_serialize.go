package api

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
