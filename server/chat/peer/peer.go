package peer

import (
	"errors"
	"fmt"
	"gowschat/server/api"
	"gowschat/server/auth/user"
	"gowschat/server/chat/room"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type (
	ChatPeer struct {
		chat        api.IChatServer
		conn        *websocket.Conn
		u           user.ChatUser
		parser      api.IParser
		serializer  api.ISerializer
		outgoing    chan api.IMessage
		peerId      string
		rooms       map[*room.ChatRoom]bool
		PeerType    api.PeerType
		status      api.PeerStatus
		registerred bool
	}
)

func NewChatPeer(chatServer api.IChatServer, peerType api.PeerType, con *websocket.Conn, u user.ChatUser) (api.IChatPeer, error) {
	return &ChatPeer{
		chat:        chatServer,
		conn:        con,
		outgoing:    make(chan api.IMessage),
		peerId:      uuid.NewString(),
		PeerType:    peerType,
		rooms:       map[*room.ChatRoom]bool{},
		status:      api.PeerStatusOnline,
		registerred: false,
		u:           u,
	}, nil
}

func (p *ChatPeer) Close() {
	p.status = api.PeerStatusOffline
	p.conn.Close()
}

func (p *ChatPeer) IsOnline() bool {
	return p.status == api.PeerStatusOnline
}

func (p *ChatPeer) TakeMessage(m api.IMessage) {
	p.outgoing <- m
}

func (p *ChatPeer) isRegisterred() bool {
	return p.registerred
}

func (p *ChatPeer) GetType() api.PeerType {
	return p.PeerType
}

func (p *ChatPeer) GetUserName() string {
	return p.u.GetUsername()
}

func GetPeerType(peerTypeString string) (api.PeerType, error) {
	if len(peerTypeString) == 0 {
		return api.PeerTypeUnknown, api.ErrUnsupporterPeerType
	}

	peerType := api.PeerType(peerTypeString)
	switch peerType {
	case api.PeerTypeJson:
		return api.PeerTypeJson, nil
	case api.PeerTypeWeb:
		return api.PeerTypeWeb, nil
	default:
		return api.PeerTypeUnknown, api.ErrUnsupporterPeerType
	}
}

func (p *ChatPeer) createEvent(m api.IMessageSerializable) (api.IEvent, error) {
	messagedata, err := m.Serialize()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if p.PeerType == api.PeerTypeJson {
		return m.CreateEvent(messagedata)
	}
	if p.PeerType == api.PeerTypeProto {
		return nil, errors.Join(api.ErrUnsupporterPeerType, fmt.Errorf("%v", p.PeerType))
	}
	return nil, api.ErrUnsupporterPeerType
}
