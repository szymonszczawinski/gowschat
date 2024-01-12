package chat

import (
	"errors"
	"gowschat/server/api"
	"log"
)

func HandlerChatRegister(e api.Event, p *ChatPeer) error {
	if message, err := e.ParseMessage(); err != nil {
		log.Println("ERROR", err)
	} else {
		if messageRegister, ok := message.(api.MessageRegister); ok {
			if err := p.server.RegisterPeer(p, messageRegister.GetEmail(), messageRegister.GetPassword()); err != nil {
				return err
			}
		} else {
			return api.ErrUnsupportedMessageType
		}
	}
	return nil
}

func HandlerMessageIn(event api.Event, p *ChatPeer) error {
	if message, err := event.ParseMessage(); err != nil {
		log.Println("ERROR", err)
	} else {
		// Broadcast to all other Clients
		p.server.BroadcastMessage(message)
	}
	return nil
}

func HandlerCreateRoom(e api.Event, p *ChatPeer) error {
	if message, parseErr := e.ParseMessage(); parseErr != nil {
		log.Println("ERROR", parseErr)
	} else {
		if createRoomMessage, ok := message.(api.MessageCreateRoom); ok {
			room, appErr := p.server.createRoom(createRoomMessage.GetRoomName(), p)
			if appErr != nil {
				return appErr
			}
			roomMessage, err := p.createMessageRoom(room)
			if err != nil {
				return err
			}
			p.outgoing <- roomMessage

		} else {
			return api.ErrUnsupportedMessageType
		}
	}
	return nil
}

func HandlerJoinRoom(e api.Event, p *ChatPeer) error {
	if message, parseErr := e.ParseMessage(); parseErr != nil {
		log.Println("ERROR", parseErr)
	} else {
		if joinRoomMessage, ok := message.(api.MessageJoinRoom); ok {
			room, appErr := p.server.joinRoom(joinRoomMessage.GetRoomName(), p)
			if appErr != nil {
				return appErr
			}
			message, err := p.createMessageRoom(room)
			if err != nil {
				return err
			}
			p.outgoing <- message

		} else {
			return api.ErrUnsupportedMessageType
		}
	}

	return nil
}

func HandlerGetRoom(e api.Event, p *ChatPeer) error {
	if message, parseErr := e.ParseMessage(); parseErr != nil {
		log.Println("ERROR", parseErr)
	} else {
		if getRoomMessage, ok := message.(api.MessageGetRoom); ok {
			room, appErr := p.server.getRoom(getRoomMessage.GetRoomName())
			if appErr != nil {
				return appErr
			}
			message, appErr := p.createMessageRoom(room)
			if appErr != nil {
				return appErr
			}
			p.outgoing <- message

		} else {
			return api.ErrUnsupportedMessageType
		}
	}

	return nil
}

func HandlerLeaveRoom(e api.Event, p *ChatPeer) error {
	// TODO: todo
	return nil
}

func HandlerRoomList(e api.Event, p *ChatPeer) error {
	message, err := p.createMessageRoomList()
	if err != nil {
		return err
	}
	p.outgoing <- message
	return nil
}

func HandlerError(err error, e api.Event, p *ChatPeer) {
	log.Printf("ERROR :: handle event %v -> %v", e.GetType(), err.Error())
	var appErr api.ChatError
	if errors.As(err, &appErr) {
		errorMessage, e := p.createMessageError(err)
		if e != nil {
			log.Printf("ERROR :: handle application error %v -> %v", appErr, err)
		}

		p.outgoing <- errorMessage
	}
}
