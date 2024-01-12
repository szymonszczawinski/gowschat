package server

import (
	"gowschat/server/chat"
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

func ServeWs(chatServer *chat.ChatServer, c *gin.Context) {
	log.Println("New connection")
	// Begin by upgrading the HTTP request
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("ERROR ::", err)
		return
	}
	peer := chat.NewChatPeer(chatServer, conn)
	chatServer.ConnectPeer(peer)
	go peer.ReadMessages()
	go peer.WriteMessages()
}

func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}
