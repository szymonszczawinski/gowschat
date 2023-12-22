package server

import (
	"log"
	"sync"
)

type ChatServer struct {
	handlers       map[string]EventHandler
	connectedPeers map[*ChatPeer]bool
	chatUsers      map[UserCredentials]*ChatUser
	rooms          map[*ChatRoom]bool
	muPeers        sync.RWMutex
	muUsers        sync.RWMutex
	muRooms        sync.RWMutex
}

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
		return ErrUserAlreadyRegisterred
	}
	chatUser := NewChatUser(p, credentials)
	chat.chatUsers[credentials] = chatUser
	return nil
}

func (chat *ChatServer) createRoom(name string, creator *ChatPeer) (*ChatRoom, error) {
	chat.muRooms.Lock()
	defer chat.muRooms.Unlock()
	for room := range chat.rooms {
		if room.name == name {
			return nil, ErrRoomExist
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
	return nil, ErrRoomNotExist
}

func (chat *ChatServer) getRoom(name string) (*ChatRoom, error) {
	chat.muRooms.Lock()
	defer chat.muRooms.Unlock()
	for room := range chat.rooms {
		if room.name == name {
			return room, nil
		}
	}
	return nil, ErrRoomNotExist
}

func (chat *ChatServer) setupHandlers() {
	chat.handlers[EventMessageIn] = HandlerMessageIn
	chat.handlers[EventRoomListGet] = HandlerRoomList
	chat.handlers[EventRoomCreate] = HandlerCreateRoom
	chat.handlers[EventRoomJoin] = HandlerJoinRoom
	chat.handlers[EventRoomGet] = HandlerGetRoom
	chat.handlers[EventChatRegister] = HandlerChatRegister
}

func (chat *ChatServer) routeEvent(e Event, p *ChatPeer) error {
	if handler, ok := chat.handlers[e.GetType()]; ok {
		if p.isRegisterred() || e.GetType() == EventChatRegister {
			if err := handler(e, p); err != nil {
				HandlerError(err, e, p)
			}
		}
		HandlerError(ErrUserNotRegisterred, e, p)
	} else {
		return ErrEventNotSupported
	}
	return nil
}
