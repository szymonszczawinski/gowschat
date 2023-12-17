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
	rooms    map[*ChatRoom]bool
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
	chat.handlers[EventInMessage] = InMessageHandler
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

func InMessageHandler(event Event, p *ChatPeer) error {
	if message, err := parseMessage(event, p.peerType); err != nil {
		log.Println("ERROR", err)
	} else {
		log.Println("message handled", message)
		outMessage := message.GenerateOutMessage()
		// Broadcast to all other Clients
		for peer := range p.server.peers {
			// if p.peerId != peer.peerId {
			peer.outgoing <- outMessage
			// }
		}
	}
	return nil
}

func CreateRoomHandler(e Event, p *ChatPeer) error {
	// TODO: todo
	return nil
}

func JoinRoomHandler(e Event, p *ChatPeer) error {
	// TODO: todo
	return nil
}

func LeaveRoomHandler(e Event, p *ChatPeer) error {
	// TODO: todo
	return nil
}
