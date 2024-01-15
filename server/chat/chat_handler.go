package chat

import (
	"gowschat/server/api"
	"gowschat/server/chat/auth"
	"gowschat/server/chat/peer"
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

func (cs *ChatServer) Login(username, password string) (auth.OTP, error) {
	otp, err := cs.authenticator.Login(username, password)

	return otp, err
}

func (cs *ChatServer) ServeWs(c *gin.Context) {
	log.Println("New connection from:", c.Request.RemoteAddr, c.Request.URL)
	otp := c.Query("otp")
	// TODO: Add authorisation based on otp parame from connection URL
	if otp == "" {
		// Tell the user its not authorized
		log.Println("ServeWs otp is empty")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Verify OTP is existing
	ok, u := cs.authenticator.VerifyOTP(otp)
	if !ok {
		log.Println("Verify OTP FAILED")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	peerTypeParam := c.Query(api.WSPeerType)
	peerType, err := peer.GetPeerType(peerTypeParam)
	if err != nil {
		log.Println("ERROR ::", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// Begin by upgrading the HTTP request
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("ERROR ::", err)
		return
	}
	log.Println("New WS connection REMOTE:", conn.RemoteAddr())
	peer, err := peer.NewChatPeer(cs, peerType, conn, u)
	if err != nil {
		log.Println("ERROR ::", err)
		conn.Close()
		return
	}
	log.Println("Peer connected:", peer)
	cs.ConnectPeer(peer)
	go peer.ReadMessages()
	go peer.WriteMessages()
}
