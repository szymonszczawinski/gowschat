package server

type (
	IParser interface {
		ParseEvent(pyload []byte) (Event, error)
		ParseMessage()
	}

	ISerializer interface{}
)
