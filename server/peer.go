package server

import (
	"log"

	"github.com/gorilla/websocket"
)

type ChatPeer struct {
	server *ChatServer
	con    *websocket.Conn
}

func NewChatPeer(chatServer *ChatServer, con *websocket.Conn) *ChatPeer {
	return &ChatPeer{
		server: chatServer,
		con:    con,
	}
}

func (p *ChatPeer) readMessages() {
	defer func() {
		p.server.RemovePeer(p)
	}()
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
	}
}
