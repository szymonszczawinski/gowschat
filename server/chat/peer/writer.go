package peer

import (
	"bytes"
	"context"
	"gowschat/server/api"
	"gowschat/server/chat/messages"
	"gowschat/server/chat/parser"
	"gowschat/server/view"
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

		case <-ticker.C:
			log.Println("ping")
			// Send the Ping
			if err := p.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("ERROR :: PING", err)
				return // return to break this goroutine triggeing cleanup
			}
		}
	}
}

func (p *ChatPeer) sendCloseMessage() error {
	return p.conn.WriteMessage(websocket.CloseMessage, nil)
}

func (p *ChatPeer) writeMessage(m api.IMessage) error {
	if p.PeerType == api.PeerTypeJson {
		event := messages.NewEvent(m)
		data, serialisationError := parser.SerializeEvent(event)
		if serialisationError != nil {
			return serialisationError
		}
		return p.conn.WriteMessage(websocket.TextMessage, data)
	}
	if p.PeerType == api.PeerTypeWeb {
		log.Println("writeMessage WEB")
		if imMessage, isIM := m.(messages.MessageIM); isIM {
			component := view.Message(imMessage.GetMessage())
			buffer := &bytes.Buffer{}
			component.Render(context.Background(), buffer)
			err := p.conn.WriteMessage(websocket.TextMessage, buffer.Bytes())
			if err != nil {
				log.Println("writeMessage error WEB", err)
			}
		}
	}
	if p.PeerType == api.PeerTypeProto {
		return api.ErrUnsupporterPeerType
	}
	return api.ErrUnsupporterPeerType
}

func (p *ChatPeer) writeError(err error) error {
	if p.PeerType == api.PeerTypeJson {
		event := messages.NewEvent(messages.NewMessageError(err))
		data, serialisationError := parser.SerializeEvent(event)
		if serialisationError != nil {
			return serialisationError
		}
		return p.conn.WriteMessage(websocket.TextMessage, data)
	}
	if p.PeerType == api.PeerTypeProto {
		return api.ErrUnsupporterPeerType
	}
	return api.ErrUnsupporterPeerType
}
