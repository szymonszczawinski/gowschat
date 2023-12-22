package server

import (
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
	PeerStatus string
	PeerType   string
)

const (
	PeerTypeJson  = "json"
	PeerTypeProto = "proto"
	PeerTypeUnset = "unset"

	PeerStatusOnline  = "online"
	PeerStatusOffline = "offline"
)

type (
	ChatPeer struct {
		server      *ChatServer
		con         *websocket.Conn
		parser      IParser
		serializer  ISerializer
		outgoing    chan MessageSerializable
		peerId      string
		rooms       map[*ChatRoom]bool
		peerType    PeerType
		status      PeerStatus
		registerred bool
	}
	UserCredentials struct {
		email    string
		password string
	}
	ChatUser struct {
		*ChatPeer
		UserCredentials
	}
)

func NewChatPeer(chatServer *ChatServer, con *websocket.Conn) *ChatPeer {
	return &ChatPeer{
		server:      chatServer,
		con:         con,
		outgoing:    make(chan MessageSerializable),
		peerId:      uuid.NewString(),
		peerType:    PeerTypeUnset,
		rooms:       map[*ChatRoom]bool{},
		status:      PeerStatusOnline,
		registerred: false,
	}
}

func NewUserCredentials(email, password string) UserCredentials {
	return UserCredentials{
		email:    email,
		password: password,
	}
}

func NewChatUser(p *ChatPeer, c UserCredentials) *ChatUser {
	return &ChatUser{
		UserCredentials: c,
		ChatPeer:        p,
	}
}

func (p *ChatPeer) readMessages() {
	defer func() {
		p.server.DisconnectPeer(p)
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
		// ReadMessage is used to read the next message in queue
		// in the connection
		messageType, payload, err := p.readMessage()
		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but not simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ERROR :: reading message: %v", err)
			}
			p.status = PeerStatusOffline
			break // Break the loop to close conn & Cleanup
		}

		event, err := parseEvent(messageType, payload)
		if err != nil {
			log.Println("ERROR :: parseEvent:", err)
		} else {
			// Route the Event
			if err := p.server.routeEvent(event, p); err != nil {
				log.Println("ERROR :: routeEvent", err)
			}
		}
	}
}

func (p *ChatPeer) writeMessages() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		p.server.DisconnectPeer(p)
		ticker.Stop()
	}()
	for {
		select {
		case messageOut, ok := <-p.outgoing:
			if !ok {
				if err := p.close(); err != nil {
					log.Println("ERROR :: connection closed", err)
				}
				return
			}
			event, err := createEvent(messageOut, p.peerType)
			if err != nil {
				log.Println("ERROR :: createEvent:", err)
				continue
			}
			data, err := event.Serialize()
			if err != nil {
				log.Println("ERROR :: Event.Serialize:", err)
				continue
			}

			if err := p.writeMessage(data); err != nil {
				log.Println("ERROR :: writeMessage:", err)
			}
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
	if p.peerType == PeerTypeUnset {
		switch mesageType {
		case websocket.TextMessage:
			p.peerType = PeerTypeJson
		case websocket.BinaryMessage:
			p.peerType = PeerTypeProto
		}
	}
}

func (p *ChatPeer) writeMessage(data []byte) error {
	if p.peerType == PeerTypeJson {
		return p.con.WriteMessage(websocket.TextMessage, data)
	}
	if p.peerType == PeerTypeProto {
		return p.con.WriteMessage(websocket.BinaryMessage, data)
	}
	return ErrUnsupporterPeerType
}

func (p *ChatPeer) readMessage() (int, []byte, error) {
	messageType, payload, err := p.con.ReadMessage()
	p.initPeer(messageType)
	return messageType, payload, err
}

func (p *ChatPeer) close() error {
	return p.con.WriteMessage(websocket.CloseMessage, nil)
}

func (p *ChatPeer) writeError(err error) error {
	if p.peerType == PeerTypeJson {
		return p.con.WriteMessage(websocket.TextMessage, []byte(err.Error()))
	}
	if p.peerType == PeerTypeProto {
		return p.con.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
	}
	return ErrUnsupporterPeerType
}
