package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

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

type ChatServer struct {
	peers   map[*ChatPeer]bool
	muPeers sync.RWMutex
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		peers: map[*ChatPeer]bool{},
	}
}

func (chat *ChatServer) AddPeer(peer *ChatPeer) {
	chat.muPeers.Lock()
	defer chat.muPeers.Unlock()
	chat.peers[peer] = true
}

func (chat *ChatServer) run() {
}

func serveWs(chatServer *ChatServer, c *gin.Context) {
	log.Println("New connection")
	// Begin by upgrading the HTTP request
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	peer := NewChatPeer(chatServer, conn)
	chatServer.AddPeer(peer)
	// We wont do anything yet so close connection again
	conn.Close()
}

func home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}
