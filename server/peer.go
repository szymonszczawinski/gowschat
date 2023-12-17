package server

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	pongWait               = 10 * time.Second
	pingInterval           = (pongWait * 9) / 10
	ErrUnsupporterPeerType = errors.New("unsupported peer type")
)

type PeerType string

const (
	PeerTypeJson  = "json"
	PeerTypeProto = "proto"
	PeerTypeUnset = "unset"
)

type ChatPeer struct {
	server   *ChatServer
	con      *websocket.Conn
	outgoing chan OutMessage
	peerId   string
	peerType PeerType
}

func NewChatPeer(chatServer *ChatServer, con *websocket.Conn) *ChatPeer {
	return &ChatPeer{
		server:   chatServer,
		con:      con,
		outgoing: make(chan OutMessage),
		peerId:   uuid.NewString(),
		peerType: PeerTypeUnset,
	}
}

func (p *ChatPeer) readMessages() {
	defer func() {
		p.server.RemovePeer(p)
	}()

	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := p.con.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
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
				log.Printf("error reading message: %v", err)
			}
			break // Break the loop to close conn & Cleanup
		}
		event, err := parseEvent(messageType, payload)
		if err != nil {
			log.Println("Error handeling Message: ", err)
		} else {
			// Route the Event
			if err := p.server.routeEvent(event, p); err != nil {
				log.Println("Error handeling Message: ", err)
			}
		}
	}
}

func (p *ChatPeer) writeMessages() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		p.server.RemovePeer(p)
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-p.outgoing:
			if !ok {
				if err := p.close(); err != nil {
					log.Println("connection closed", err)
				}
				return
			}

			messagedata, err := serializeOutMessage(message, p.peerType)
			if err != nil {
				log.Println(err)
			}

			event, err := wrapOutMessage(messagedata, p.peerType)
			if err != nil {
				log.Println(err)
			}
			data, err := event.Serialize()
			if err != nil {
				log.Println(err)
			}

			if err := p.writeMessage(data); err != nil {
				log.Println(err)
			}
			log.Println("message sent")
		case <-ticker.C:
			log.Println("ping")
			// Send the Ping
			if err := p.con.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("writemsg: ", err)
				return // return to break this goroutine triggeing cleanup
			}
		}
	}
}

// pongHandler is used to handle PongMessages for the Client
func (p *ChatPeer) pongHandler(pongMsg string) error {
	// Current time + Pong Wait time
	log.Println("pong", pongMsg)
	return p.con.SetReadDeadline(time.Now().Add(pongWait))
}

func (p *ChatPeer) sendInitalMessage() {
	event := generateInitialEvent(*p)
	data, err := event.Serialize()
	if err != nil {
		log.Println(err)
		return // closes the connection, should we really
	}
	if err := p.con.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println(err)
	}
	log.Println("message sent")
}

func (p *ChatPeer) initPeerType(mesageType int) {
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
	p.initPeerType(messageType)
	return messageType, payload, err
}

func (p *ChatPeer) close() error {
	return p.con.WriteMessage(websocket.CloseMessage, nil)
}
