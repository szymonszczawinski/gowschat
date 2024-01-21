package api

type (
	PeerStatus string
	PeerType   string
)

const (
	WSPeerType               = "peerType"
	PeerTypeJson    PeerType = "json"
	PeerTypeProto   PeerType = "proto"
	PeerTypeWeb     PeerType = "web"
	PeerTypeUnset   PeerType = "unset"
	PeerTypeUnknown PeerType = "unnkonw"

	PeerStatusOnline  PeerStatus = "online"
	PeerStatusOffline PeerStatus = "offline"
)

type (
	IChatServer interface {
		ConnectPeer(p IChatPeer)
		DisconnectPeer(p IChatPeer)
		RouteEvent(e IEvent, p IChatPeer) error
		BroadcastMessage(m IMessage) error
	}
	IChatPeer interface {
		TakeMessage(m IMessage)
		Close()
		IsOnline() bool
		ReadMessages()
		WriteMessages()
		GetType() PeerType
		GetUserName() string
	}
	IChatUser        interface{}
	IUserCredentials interface{}
)
