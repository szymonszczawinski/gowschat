package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeWs(chatServer *ChatServer, c *gin.Context) {
	log.Println("New connection")
	// Begin by upgrading the HTTP request
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	peer := NewChatPeer(chatServer, conn)
	chatServer.AddPeer(peer)
	go peer.readMessages()
	go peer.writeMessages()
}

func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}
