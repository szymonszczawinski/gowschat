package peer

import (
	"gowschat/server/api"
	"gowschat/server/chat"
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
		messageType, payload, err := p.con.ReadMessage()
		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but not simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ERROR :: reading message: %v", err)
			}
			p.status = api.PeerStatusOffline
			break // Break the loop to close conn & Cleanup
		}
		event, err := chat.ParseIncomingMessage(p, messageType, payload)
		if err != nil {
			log.Println("ERROR :: parseEvent:", err)
			p.writeError(err)
		} else {
			// Route the Event
			if err := p.chat.RouteEvent(event, p); err != nil {
				log.Println("ERROR :: routeEvent", err)
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
