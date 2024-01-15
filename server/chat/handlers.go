package chat

import (
	"gowschat/server/api"
)

// import (
//
//	"errors"
//	"gowschat/server/api"
//	"log"
//
// )
//
//	func HandlerChatRegister(chat *ChatServer, e api.Event, p api.IChatPeer) error {
//		if message, err := e.ParseMessage(); err != nil {
//			log.Println("ERROR", err)
//		} else {
//			if messageRegister, ok := message.(api.IMessageRegister); ok {
//				if err := chat.RegisterPeer(p, messageRegister.GetEmail(), messageRegister.GetPassword()); err != nil {
//					return err
//				}
//			} else {
//				return api.ErrUnsupportedMessageType
//			}
//		}
//		return nil
//	}
func HandlerMessageIM(chat *ChatServer, event api.IEvent, p api.IChatPeer) error {
	// Broadcast to all other Clients
	chat.BroadcastMessage(event.GetMessage())
	return nil
}

//
// func HandlerCreateRoom(chat *ChatServer, e api.Event, p api.IChatPeer) error {
// 	if message, parseErr := e.ParseMessage(); parseErr != nil {
// 		log.Println("ERROR", parseErr)
// 	} else {
// 		if createRoomMessage, ok := message.(api.IMessageCreateRoom); ok {
// 			room, appErr := chat.createRoom(createRoomMessage.GetRoomName(), p)
// 			if appErr != nil {
// 				return appErr
// 			}
// 			roomMessage, err := chat.createMessageRoom(room)
// 			if err != nil {
// 				return err
// 			}
// 			p.TakeMessage(roomMessage)
// 			// p.outgoing <- roomMessage
//
// 		} else {
// 			return api.ErrUnsupportedMessageType
// 		}
// 	}
// 	return nil
// }
//
// func HandlerJoinRoom(chat *ChatServer, e api.Event, p api.IChatPeer) error {
// 	if message, parseErr := e.ParseMessage(); parseErr != nil {
// 		log.Println("ERROR", parseErr)
// 	} else {
// 		if joinRoomMessage, ok := message.(api.IMessageJoinRoom); ok {
// 			room, appErr := chat.joinRoom(joinRoomMessage.GetRoomName(), p)
// 			if appErr != nil {
// 				return appErr
// 			}
// 			message, err := chat.createMessageRoom(room)
// 			if err != nil {
// 				return err
// 			}
// 			// p.outgoing <- message
// 			p.TakeMessage(message)
//
// 		} else {
// 			return api.ErrUnsupportedMessageType
// 		}
// 	}
//
// 	return nil
// }
//
// func HandlerGetRoom(chat *ChatServer, e api.Event, p api.IChatPeer) error {
// 	if message, parseErr := e.ParseMessage(); parseErr != nil {
// 		log.Println("ERROR", parseErr)
// 	} else {
// 		if getRoomMessage, ok := message.(api.IMessageGetRoom); ok {
// 			room, appErr := chat.getRoom(getRoomMessage.GetRoomName())
// 			if appErr != nil {
// 				return appErr
// 			}
// 			message, appErr := chat.createMessageRoom(room)
// 			if appErr != nil {
// 				return appErr
// 			}
// 			// p.outgoing <- message
// 			p.TakeMessage(message)
//
// 		} else {
// 			return api.ErrUnsupportedMessageType
// 		}
// 	}
//
// 	return nil
// }
//
// func HandlerLeaveRoom(chat *ChatServer, e api.Event, p api.IChatPeer) error {
// 	// TODO: todo
// 	return nil
// }
//
// func HandlerRoomList(chat *ChatServer, e api.Event, p api.IChatPeer) error {
// 	message, err := chat.createMessageRoomList()
// 	if err != nil {
// 		return err
// 	}
// 	// p.outgoing <- message
// 	p.TakeMessage(message)
// 	return nil
// }
//
// func HandlerError(chat *ChatServer, err error, e api.Event, p api.IChatPeer) {
// 	log.Printf("ERROR :: handle event %v -> %v", e.GetType(), err.Error())
// 	var appErr api.ChatError
// 	if errors.As(err, &appErr) {
// 		errorMessage, e := chat.createMessageError(err)
// 		if e != nil {
// 			log.Printf("ERROR :: handle application error %v -> %v", appErr, err)
// 		}
//
// 		// p.outgoing <- errorMessage
// 		p.TakeMessage(errorMessage)
// 	}
// }
