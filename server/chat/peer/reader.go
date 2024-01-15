package peer

import (
	"errors"
	"gowschat/server/api"
	"gowschat/server/chat/parser"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

func (p *ChatPeer) ReadMessages() {
	log.Println("ReadMessages")
	defer func() {
		p.chat.DisconnectPeer(p)
	}()
	setConnctionConfig(p.conn)
	// Loop Forever
	for {
		// ReadMessage is used to read the next message in queue
		// in the connection
		messageType, payload, err := p.conn.ReadMessage()
		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but not simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("ERROR :: reading message", err)
			}
			p.status = api.PeerStatusOffline
			break // Break the loop to close conn & Cleanup
		}
		event, err := parser.ParseIncomingMessage(p, messageType, payload)
		if err != nil {
			log.Println("ERROR :: parseEvent:", err)
			if !errors.Is(err, api.ErrEmptyMessage) {
				p.writeError(err)
			}
		} else {
			// Route the Event
			if err := p.chat.RouteEvent(event, p); err != nil {
				log.Println("ERROR :: routeEvent", err)
			}
		}
	}
}

// pongHandler is used to handle PongMessages for the Client
func pongHandler(conn *websocket.Conn, pongMsg string) func(string) error {
	return func(string) error {
		// Current time + Pong Wait time
		log.Println("pong", pongMsg)
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	}
}

func setConnctionConfig(conn *websocket.Conn) {
	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println("ERROR :: SetReadDeadline:", err)
		return
	}
	// Configure how to handle Pong responses
	conn.SetPongHandler(pongHandler(conn, "OK"))
	conn.SetReadLimit(512)
}
