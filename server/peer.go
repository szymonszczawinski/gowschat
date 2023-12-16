package server

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ChatPeer struct {
	server   *ChatServer
	con      *websocket.Conn
	outgoing chan Event
}

func NewChatPeer(chatServer *ChatServer, con *websocket.Conn) *ChatPeer {
	return &ChatPeer{
		server:   chatServer,
		con:      con,
		outgoing: make(chan Event),
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
	// Configure how to handle Pong responses
	p.con.SetPongHandler(p.pongHandler)

	// Loop Forever
	for {
		// ReadMessage is used to read the next message in queue
		// in the connection
		messageType, payload, err := p.con.ReadMessage()
		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but not simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break // Break the loop to close conn & Cleanup
		}
		log.Println("MessageType: ", messageType)
		log.Println("Payload: ", string(payload))
		request, err := parseEvent(messageType, payload)
		if err != nil {
			log.Println("Error handeling Message: ", err)
		}
		// Route the Event
		if err := p.server.routeEvent(request, p); err != nil {
			log.Println("Error handeling Message: ", err)
		}
	}
}

// TODO: xxx
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
				if err := p.con.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed", err)
				}
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return // closes the connection, should we really
			}
			if err := p.con.WriteMessage(websocket.TextMessage, data); err != nil {
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
