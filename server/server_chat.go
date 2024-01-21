package server

import (
	"gowschat/server/api"
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

func serveWs(c *gin.Context, s *server) {
	log.Println("New connection from:", c.Request.RemoteAddr, c.Request.URL)
	otp := c.Query(OTP_KEY)
	if otp == "" {
		// Tell the user its not authorized
		log.Println("ERROR :: otp is empty")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Verify OTP is existing
	ok, user := s.authenticator.VerifyOTP(otp)
	if !ok {
		log.Println("ERROR :: Verify OTP FAILED")
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
	s.chat.NewPeer(conn, peerType, user)
}
