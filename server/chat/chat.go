package chat

import (
	"gowschat/server/api"
	"log"
	"sync"
)

type (
	ChatServer struct {
		handlers       map[string]EventHandler
		connectedPeers map[*ChatPeer]bool
		chatUsers      map[UserCredentials]*ChatUser
		rooms          map[*ChatRoom]bool
		muPeers        sync.RWMutex
		muUsers        sync.RWMutex
		muRooms        sync.RWMutex
	}

	EventHandler func(event api.Event, p *ChatPeer) error
)

func NewChatServer() *ChatServer {
	server := &ChatServer{
		connectedPeers: map[*ChatPeer]bool{},
		chatUsers:      map[UserCredentials]*ChatUser{},
		handlers:       map[string]EventHandler{},
		rooms:          map[*ChatRoom]bool{},
	}
	server.setupHandlers()
	return server
}

func (chat *ChatServer) ConnectPeer(peer *ChatPeer) {
	chat.muPeers.Lock()
	defer chat.muPeers.Unlock()
	chat.connectedPeers[peer] = true
}

func (chat *ChatServer) DisconnectPeer(peer *ChatPeer) {
	chat.muPeers.Lock()
	defer chat.muPeers.Unlock()
	if _, exist := chat.connectedPeers[peer]; exist {
		peer.con.Close()
		delete(chat.connectedPeers, peer)
	}
}

func (chat *ChatServer) RegisterPeer(p *ChatPeer, email, password string) error {
	chat.muUsers.Lock()
	defer chat.muUsers.Unlock()
	credentials := NewUserCredentials(email, password)
	if _, exist := chat.chatUsers[credentials]; exist {
		return api.ErrUserAlreadyRegisterred
	}
	chatUser := NewChatUser(p, credentials)
	chat.chatUsers[credentials] = chatUser
	return nil
}

func (chat *ChatServer) RouteEvent(e api.Event, p *ChatPeer) error {
	if handler, ok := chat.handlers[e.GetType()]; ok {
		// if p.isRegisterred() || e.GetType() == api.EventChatRegister {
		if err := handler(e, p); err != nil {
			HandlerError(err, e, p)
		}
		// }
		// HandlerError(api.ErrUserNotRegisterred, e, p)
	} else {
		return api.ErrEventNotSupported
	}
	return nil
}

func (chat *ChatServer) BroadcastMessage(m api.Message) error {
	for peer := range chat.connectedPeers {
		messageOut, err := peer.createMessageOut(m)
		if err != nil {
			return err
		}
		// FIXME:: re-enable same peer check
		// if p.peerId != peer.peerId {
		if peer.status == PeerStatusOnline {
			peer.outgoing <- messageOut
		}
	}
	return nil
}

func (chat *ChatServer) createRoom(name string, creator *ChatPeer) (*ChatRoom, error) {
	chat.muRooms.Lock()
	defer chat.muRooms.Unlock()
	for room := range chat.rooms {
		if room.name == name {
			return nil, api.ErrRoomExist
		}
	}
	newRoom := NewChatRoom(name, creator)
	chat.rooms[newRoom] = true
	log.Printf("room created %v\n", name)
	return newRoom, nil
}

func (chat *ChatServer) joinRoom(name string, joiner *ChatPeer) (*ChatRoom, error) {
	chat.muRooms.Lock()
	defer chat.muRooms.Unlock()
	for room := range chat.rooms {
		if room.name == name {
			if err := room.join(joiner); err != nil {
				return nil, err
			}
			return room, nil
		}
	}
	return nil, api.ErrRoomNotExist
}

func (chat *ChatServer) getRoom(name string) (*ChatRoom, error) {
	chat.muRooms.Lock()
	defer chat.muRooms.Unlock()
	for room := range chat.rooms {
		if room.name == name {
			return room, nil
		}
	}
	return nil, api.ErrRoomNotExist
}

func (chat *ChatServer) setupHandlers() {
	chat.handlers[api.EventMessageIn] = HandlerMessageIn
	chat.handlers[api.EventRoomListGet] = HandlerRoomList
	chat.handlers[api.EventRoomCreate] = HandlerCreateRoom
	chat.handlers[api.EventRoomJoin] = HandlerJoinRoom
	chat.handlers[api.EventRoomGet] = HandlerGetRoom
	chat.handlers[api.EventChatRegister] = HandlerChatRegister
}
