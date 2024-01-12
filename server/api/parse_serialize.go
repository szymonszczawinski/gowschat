package api

type (
	IParser interface {
		ParseEvent(messageType int, pyload []byte) (Event, error)
	}

	ISerializer interface{}
)
