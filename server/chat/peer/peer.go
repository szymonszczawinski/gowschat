package peer

import (
	"errors"
	"fmt"
	"gowschat/server/api"
	"gowschat/server/chat"
	"gowschat/server/chat/room"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type (
	ChatPeer struct {
		chat        *chat.ChatServer
		con         *websocket.Conn
		parser      api.IParser
		serializer  api.ISerializer
		outgoing    chan api.IMessage
		peerId      string
		rooms       map[*room.ChatRoom]bool
		PeerType    api.PeerType
		status      api.PeerStatus
		registerred bool
	}
)

func NewChatPeer(chatServer *chat.ChatServer, con *websocket.Conn) api.IChatPeer {
	return &ChatPeer{
		chat:        chatServer,
		con:         con,
		outgoing:    make(chan api.IMessage),
		peerId:      uuid.NewString(),
		PeerType:    api.PeerTypeUnset,
		rooms:       map[*room.ChatRoom]bool{},
		status:      api.PeerStatusOnline,
		registerred: false,
	}
}

func (p *ChatPeer) Close() {
	p.con.Close()
}

func (p *ChatPeer) IsOnline() bool {
	return p.status == api.PeerStatusOnline
}

func (p *ChatPeer) TakeMessage(m api.IMessage) {
	p.outgoing <- m
}

func (p *ChatPeer) ReadMessages() {
	log.Println("ReadMessages")
	defer func() {
		p.chat.DisconnectPeer(p)
	}()

	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := p.con.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println("ERROR :: SetReadDeadline:", err)
		return
	}
	p.con.SetReadLimit(512)
	// Configure how to handle Pong responses
	p.con.SetPongHandler(p.pongHandler)

	// Loop Forever
	for {

		log.Println("ReadMessages FOR LOOP")
		// ReadMessage is used to read the next message in queue
		// in the connection
		messageType, payload, err := p.readMessage()

		log.Println("ReadMessages MESSAGE READ")
		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but not simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ERROR :: reading message: %v", err)
			}
			p.status = api.PeerStatusOffline
			break // Break the loop to close conn & Cleanup
		}
		event, err := chat.ParseIncomingMessage(messageType, payload)
		log.Println(event, err)
		if err != nil {
			log.Println("ERROR :: parseEvent:", err)
		} else {
			// Route the Event
			if err := p.chat.RouteEvent(event, p); err != nil {
				log.Println("ERROR :: routeEvent", err)
			}
		}
	}
}

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

func (p *ChatPeer) isRegisterred() bool {
	return p.registerred
}

// pongHandler is used to handle PongMessages for the Client
func (p *ChatPeer) pongHandler(pongMsg string) error {
	// Current time + Pong Wait time
	log.Println("pong", pongMsg)
	return p.con.SetReadDeadline(time.Now().Add(pongWait))
}

func (p *ChatPeer) initPeer(mesageType int) {
	if p.PeerType == api.PeerTypeUnset {
		switch mesageType {
		case websocket.TextMessage:
			p.PeerType = api.PeerTypeJson
			// p.parser = &chatjson.ParserJson{}
		case websocket.BinaryMessage:
			p.PeerType = api.PeerTypeProto
		}
	}
}

func (p *ChatPeer) writeMessage(data []byte) error {
	if p.PeerType == api.PeerTypeJson {
		return p.con.WriteMessage(websocket.TextMessage, data)
	}
	if p.PeerType == api.PeerTypeProto {
		return p.con.WriteMessage(websocket.BinaryMessage, data)
	}
	return api.ErrUnsupporterPeerType
}

func (p *ChatPeer) readMessage() (int, []byte, error) {
	messageType, payload, err := p.con.ReadMessage()
	p.initPeer(messageType)
	return messageType, payload, err
}

func (p *ChatPeer) sendCloseMessage() error {
	return p.con.WriteMessage(websocket.CloseMessage, nil)
}

func (p *ChatPeer) writeError(err error) error {
	if p.PeerType == api.PeerTypeJson {
		return p.con.WriteMessage(websocket.TextMessage, []byte(err.Error()))
	}
	if p.PeerType == api.PeerTypeProto {
		return p.con.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
	}
	return api.ErrUnsupporterPeerType
}

func (p *ChatPeer) createEvent(m api.IMessageSerializable) (api.IEvent, error) {
	messagedata, err := m.Serialize()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if p.PeerType == api.PeerTypeJson {
		return m.CreateEvent(messagedata)
	}
	if p.PeerType == api.PeerTypeProto {
		return nil, errors.Join(api.ErrUnsupporterPeerType, fmt.Errorf("%v", p.PeerType))
	}
	return nil, api.ErrUnsupporterPeerType
}
