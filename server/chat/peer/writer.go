package peer

import (
	"gowschat/server/api"
	"gowschat/server/chat"
	"gowschat/server/chat/messages"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (p *ChatPeer) WriteMessages() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		p.chat.DisconnectPeer(p)
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-p.outgoing:
			if !ok {
				if err := p.sendCloseMessage(); err != nil {
					log.Println("ERROR :: connection closed", err)
				}
				return
			}
			log.Println("MESSAGE OUT ::", message)
			p.writeMessage(message)
		// event, err := p.createEvent(messageOut)
		// if err != nil {
		// 	log.Println("ERROR :: createEvent:", err)
		// 	continue
		// }
		// data, err := event.Serialize()
		// if err != nil {
		// 	log.Println("ERROR :: Event.Serialize:", err)
		// 	continue
		// }
		//
		// if err := p.writeMessage(data); err != nil {
		// 	log.Println("ERROR :: writeMessage:", err)
		// }
		case <-ticker.C:
			log.Println("ping")
			// Send the Ping
			if err := p.con.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("ERROR :: PING", err)
				return // return to break this goroutine triggeing cleanup
			}
		}
	}
}

func (p *ChatPeer) sendCloseMessage() error {
	return p.con.WriteMessage(websocket.CloseMessage, nil)
}

func (p *ChatPeer) writeMessage(m api.IMessage) error {
	if p.PeerType == api.PeerTypeJson {
		event := messages.NewEvent(api.EventMessageIM, m)
		data, serialisationError := chat.SerializeEvent(event)
		if serialisationError != nil {
			return serialisationError
		}
		return p.con.WriteMessage(websocket.TextMessage, data)
	}
	if p.PeerType == api.PeerTypeProto {
		return api.ErrUnsupporterPeerType
	}
	return api.ErrUnsupporterPeerType
}

func (p *ChatPeer) writeError(err error) error {
	if p.PeerType == api.PeerTypeJson {
		event := messages.NewEvent(api.EventError, messages.NewMessageError(err))
		data, serialisationError := chat.SerializeEvent(event)
		if serialisationError != nil {
			return serialisationError
		}
		return p.con.WriteMessage(websocket.TextMessage, data)
	}
	if p.PeerType == api.PeerTypeProto {
		return api.ErrUnsupporterPeerType
	}
	return api.ErrUnsupporterPeerType
}