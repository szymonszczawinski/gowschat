package server

import (
	"errors"
	"sync"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
	ErrRoomExist         = errors.New("room witch such name exist")
	ErrRoomNotExist      = errors.New("room does not exist")
)

type ChatServer struct {
	handlers map[string]EventHandler
	peers    map[*ChatPeer]bool
	rooms    map[*ChatRoom]bool
	muPeers  sync.RWMutex
	muRooms  sync.RWMutex
}

func NewChatServer() *ChatServer {
	server := &ChatServer{
		peers:    map[*ChatPeer]bool{},
		handlers: map[string]EventHandler{},
		rooms:    map[*ChatRoom]bool{},
	}
	server.setupHandlers()
	return server
}

func (chat *ChatServer) AddPeer(peer *ChatPeer) {
	chat.muPeers.Lock()
	defer chat.muPeers.Unlock()
	chat.peers[peer] = true
}

func (chat *ChatServer) RemovePeer(peer *ChatPeer) {
	chat.muPeers.Lock()
	defer chat.muPeers.Unlock()
	if _, exist := chat.peers[peer]; exist {
		peer.con.Close()
		delete(chat.peers, peer)
	}
}

func (chat *ChatServer) createRoom(name string, creator *ChatPeer) error {
	chat.muRooms.Lock()
	defer chat.muRooms.Unlock()
	for room := range chat.rooms {
		if room.name == name {
			return ErrRoomExist
		}
	}
	newRoom := NewChatRoom(name, creator)
	chat.rooms[newRoom] = true
	return nil
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

func (chat *ChatServer) Run() {
}

func (chat *ChatServer) setupHandlers() {
	chat.handlers[EventInMessage] = HandlerMessageIn
	chat.handlers[EventGetRoomList] = HandlerRoomList
	chat.handlers[EventCreateRoom] = HandlerCreateRoom
	chat.handlers[EventJoinRoom] = HandlerJoinRoom
	chat.handlers[EventGetRoom] = HandlerGetRoom
}

func (chat *ChatServer) routeEvent(e Event, p *ChatPeer) error {
	if handler, ok := chat.handlers[e.GetType()]; ok {
		if err := handler(e, p); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}
