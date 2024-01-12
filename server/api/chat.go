package api

type (
	PeerStatus string
	PeerType   string
)

const (
	PeerTypeJson  = "json"
	PeerTypeProto = "proto"
	PeerTypeWeb   = "web"
	PeerTypeUnset = "unset"

	PeerStatusOnline  = "online"
	PeerStatusOffline = "offline"
)

type (
	IChatServer interface {
		ConnectPeer(p IChatPeer)
		DisconnectPeer(p IChatPeer)
		RegisterPeer(p IChatPeer, email, password string) error
		RouteEvent(e IEvent, p IChatPeer) error
		BroadcastMessage(m IMessage) error
	}
	IChatPeer interface {
		TakeMessage(m IMessage)
		Close()
		IsOnline() bool
		ReadMessages()
		WriteMessages()
	}
	IChatUser        interface{}
	IUserCredentials interface{}
)
