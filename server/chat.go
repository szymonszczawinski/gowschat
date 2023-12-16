package server

import (
	"errors"
	"log"
	"sync"
)

var ErrEventNotSupported = errors.New("this event type is not supported")

type ChatServer struct {
	handlers map[string]EventHandler
	peers    map[*ChatPeer]bool
	muPeers  sync.RWMutex
}

func NewChatServer() *ChatServer {
	server := &ChatServer{
		peers:    map[*ChatPeer]bool{},
		handlers: map[string]EventHandler{},
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

func (chat *ChatServer) Run() {
}

func (chat *ChatServer) setupHandlers() {
	chat.handlers[EventSendMessage] = func(event Event, p *ChatPeer) error {
		log.Println(event)
		return nil
	}
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
