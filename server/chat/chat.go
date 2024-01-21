package chat

import (
	"gowschat/server/api"
	"gowschat/server/auth/user"
	"gowschat/server/chat/peer"
	"gowschat/server/chat/room"
	"log"

	"github.com/gorilla/websocket"
)

type (
	ChatServer struct {
		handlers       map[api.EventType]EventHandler
		connectedPeers map[api.IChatPeer]bool
		rooms          map[*room.ChatRoom]bool
		register       chan api.IChatPeer
		unregister     chan api.IChatPeer
	}
	EventHandler func(chat *ChatServer, event api.IEvent, p api.IChatPeer) error

	ClientListEvent struct {
		peer      api.IChatPeer
		EventType string
	}
)

func NewChatServer() *ChatServer {
	server := &ChatServer{
		connectedPeers: map[api.IChatPeer]bool{},
		handlers:       map[api.EventType]EventHandler{},
		rooms:          map[*room.ChatRoom]bool{},
		register:       make(chan api.IChatPeer),
		unregister:     make(chan api.IChatPeer),
	}
	server.setupHandlers()
	return server
}

// Run our websocket server, accepting various requests
func (chat *ChatServer) Run() {
	for {
		select {

		case client := <-chat.register:
			chat.registerPeer(client)

		case client := <-chat.unregister:
			chat.unregisterPeer(client)
		}
	}
}

func (chat *ChatServer) NewPeer(conn *websocket.Conn, peerType api.PeerType, user user.ChatUser) {
	peer, err := peer.NewChatPeer(chat, peerType, conn, user)
	if err != nil {
		log.Println("ERROR ::", err)
		conn.Close()
		return
	}
	log.Println("Peer connected:", peer)
	chat.ConnectPeer(peer)
	go peer.ReadMessages()
	go peer.WriteMessages()
}

func (chat *ChatServer) ConnectPeer(peer api.IChatPeer) {
	chat.register <- peer
}

func (chat *ChatServer) DisconnectPeer(peer api.IChatPeer) {
	chat.unregister <- peer
}

func (chat *ChatServer) RouteEvent(e api.IEvent, p api.IChatPeer) error {
	if handler, ok := chat.handlers[e.GetType()]; ok {
		// if p.isRegisterred() || e.GetType() == api.EventChatRegister {
		if err := handler(chat, e, p); err != nil {
			log.Fatal("ERROR ::", err)
			// HandlerError(chat, err, e, p)
		}
		// }
		// HandlerError(api.ErrUserNotRegisterred, e, p)
	} else {
		return api.ErrEventNotSupported
	}
	return nil
}

func (chat *ChatServer) BroadcastMessage(m api.IMessage) error {
	for peer := range chat.connectedPeers {
		// FIXME:: re-enable same peer check
		// if p.peerId != peer.peerId {
		if peer.IsOnline() {
			// peer.outgoing <- messageOut
			peer.TakeMessage(m)
		}
	}
	return nil
}

func (chat *ChatServer) registerPeer(peer api.IChatPeer) {
	chat.connectedPeers[peer] = true
}

func (chat *ChatServer) unregisterPeer(peer api.IChatPeer) {
	if _, exist := chat.connectedPeers[peer]; exist {
		peer.Close()
		delete(chat.connectedPeers, peer)
	}
}

//	func (chat *ChatServer) createRoom(name string, creator api.IChatPeer) (*room.ChatRoom, error) {
//		chat.muRooms.Lock()
//		defer chat.muRooms.Unlock()
//		for room := range chat.rooms {
//			if room.Name == name {
//				return nil, api.ErrRoomExist
//			}
//		}
//		newRoom := room.NewChatRoom(name, creator)
//		chat.rooms[newRoom] = true
//		log.Printf("room created %v\n", name)
//		return newRoom, nil
//	}
//
//	func (chat *ChatServer) joinRoom(name string, joiner api.IChatPeer) (*room.ChatRoom, error) {
//		chat.muRooms.Lock()
//		defer chat.muRooms.Unlock()
//		for room := range chat.rooms {
//			if room.Name == name {
//				if err := room.Join(joiner); err != nil {
//					return nil, err
//				}
//				return room, nil
//			}
//		}
//		return nil, api.ErrRoomNotExist
//	}
//
//	func (chat *ChatServer) getRoom(name string) (*room.ChatRoom, error) {
//		chat.muRooms.Lock()
//		defer chat.muRooms.Unlock()
//		for room := range chat.rooms {
//			if room.Name == name {
//				return room, nil
//			}
//		}
//		return nil, api.ErrRoomNotExist
//	}
func (chat *ChatServer) setupHandlers() {
	chat.handlers[api.EventMessageIM] = HandlerMessageIM
	// chat.handlers[api.EventRoomListGet] = HandlerRoomList
	// chat.handlers[api.EventRoomCreate] = HandlerCreateRoom
	// chat.handlers[api.EventRoomJoin] = HandlerJoinRoom
	// chat.handlers[api.EventRoomGet] = HandlerGetRoom
	// chat.handlers[api.EventChatRegister] = HandlerChatRegister
}
