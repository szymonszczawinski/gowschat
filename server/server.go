package server

import (
	"context"
	"gowschat/server/api"
	"gowschat/server/chat"
	"gowschat/server/chat/peer"
	"gowschat/server/view"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Run() {
	chatServer := chat.NewChatServer()
	// go chatServer.run()
	router := gin.Default()

	router.GET("/chat", Chat)
	router.GET("/ws/chat", func(ctx *gin.Context) {
		ServeWs(chatServer, ctx)
	})
	router.GET("/", Home)

	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	log.Println("Server exiting")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeWs(chatServer *chat.ChatServer, c *gin.Context) {
	log.Println("New connection from:", c.Request.RemoteAddr, c.Request.URL)

	// TODO: Add authorisation based on otp parame from connection URL

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
	peer, err := peer.NewChatPeer(chatServer, peerType, conn)
	if err != nil {
		log.Println("ERROR ::", err)
		conn.Close()
		return
	}
	log.Println("Peer connected:", peer)
	chatServer.ConnectPeer(peer)
	go peer.ReadMessages()
	go peer.WriteMessages()
}

func Home(c *gin.Context) {
	view.Login().Render(c.Request.Context(), c.Writer)
}

func Chat(c *gin.Context) {
	component := view.Chat("szymon")
	component.Render(c.Request.Context(), c.Writer)
}
