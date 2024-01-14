package chat

import (
	"gowschat/server/api"
	"gowschat/server/chat/messages"
	"gowschat/server/chat/room"
	"gowschat/server/chat/user"
	"log"
	"sync"
)

type (
	ChatServer struct {
		handlers       map[api.EventType]EventHandler
		connectedPeers map[api.IChatPeer]bool
		chatUsers      map[api.IUserCredentials]api.IChatUser
		rooms          map[*room.ChatRoom]bool
		muPeers        sync.RWMutex
		muUsers        sync.RWMutex
		muRooms        sync.RWMutex
	}

	EventHandler func(chat *ChatServer, event messages.Event, p api.IChatPeer) error

	ClientListEvent struct {
		peer      api.IChatPeer
		EventType string
	}
)

func NewChatServer() *ChatServer {
	server := &ChatServer{
		connectedPeers: map[api.IChatPeer]bool{},
		chatUsers:      map[api.IUserCredentials]api.IChatUser{},
		handlers:       map[api.EventType]EventHandler{},
		rooms:          map[*room.ChatRoom]bool{},
	}
	server.setupHandlers()
	return server
}

func (chat *ChatServer) ConnectPeer(peer api.IChatPeer) {
	chat.muPeers.Lock()
	defer chat.muPeers.Unlock()
	chat.connectedPeers[peer] = true
}

func (chat *ChatServer) DisconnectPeer(peer api.IChatPeer) {
	chat.muPeers.Lock()
	defer chat.muPeers.Unlock()
	if _, exist := chat.connectedPeers[peer]; exist {
		peer.Close()
		delete(chat.connectedPeers, peer)
	}
}

func (chat *ChatServer) RegisterPeer(p api.IChatPeer, email, password string) error {
	chat.muUsers.Lock()
	defer chat.muUsers.Unlock()
	credentials := user.NewUserCredentials(email, password)
	if _, exist := chat.chatUsers[credentials]; exist {
		return api.ErrUserAlreadyRegisterred
	}
	chatUser := user.NewChatUser(credentials)
	chat.chatUsers[credentials] = chatUser
	return nil
}

func (chat *ChatServer) RouteEvent(e messages.Event, p api.IChatPeer) error {
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
