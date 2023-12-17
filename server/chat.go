package server

import (
	"errors"
	"log"
	"sync"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
	ErrRoomExist         = errors.New("room witch such name exist")
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

func (chat *ChatServer) Run() {
}

func (chat *ChatServer) setupHandlers() {
	chat.handlers[EventInMessage] = InMessageHandler
	chat.handlers[EventListRooms] = ListRoomsHandler
	chat.handlers[EventCreateRoom] = CreateRoomHandler
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
		if outMessage, err := createOutMessage(message); err != nil {
			log.Println("out message not created", err)
		} else {
			// Broadcast to all other Clients
			for peer := range p.server.peers {
				// FIXME:: re-enable same peer check
				// if p.peerId != peer.peerId {
				if peer.status == PeerStatusOnline {
					peer.outgoing <- outMessage
				}
				// }
			}
		}
	}
	return nil
}

func CreateRoomHandler(e Event, p *ChatPeer) error {
	if message, err := parseMessage(e, p.peerType); err != nil {
		log.Println("ERROR", err)
	} else {
		if createRoomMessage, ok := message.(CreateRoomMessage); ok {
			if err := p.server.createRoom(createRoomMessage.GetRoomName(), p); err != nil {
				return err
			} else {
				log.Println("room created", createRoomMessage.GetRoomName())
			}
		} else {
			return ErrUnsupportedMessageType
		}
	}
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

func ListRoomsHandler(e Event, p *ChatPeer) error {
	message := createRoomListMessage(p)
	p.outgoing <- message
	return nil
}
