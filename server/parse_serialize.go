package server

import "gowschat/server/api"

type (
	IParser interface {
		ParseEvent(messageType int, pyload []byte) (api.Event, error)
	}

	ISerializer interface{}
)
