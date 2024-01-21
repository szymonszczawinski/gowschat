package server

import (
	"context"
	"gowschat/server/auth"
	"gowschat/server/chat"
	"gowschat/server/view"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func RunApp() {
	mainContext := context.Background()
	ctx, cancel := context.WithCancel(mainContext)
	defer func() {
		log.Println("Cancel 1")
		cancel()
		time.Sleep(1 * time.Second)
		log.Println("DONE")
	}()
	newServer(ctx).run()
	log.Println("exiting")
}

type server struct {
	ctx           context.Context
	chat          *chat.ChatServer
	authenticator *auth.Authenticator
	router        *gin.Engine
}

func newServer(ctx context.Context) *server {
	return &server{
		ctx:           ctx,
		chat:          chat.NewChatServer(),
		authenticator: auth.NewAuthenticator(ctx, 5*time.Second),
		router:        gin.Default(),
	}
}

func (s *server) run() {
	go s.chat.Run()
	s.configureRoutes()

	server := http.Server{
		Addr:    ":3000",
		Handler: s.router,
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

	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer func() {
		log.Println("CANCEL 2")
		cancel()
	}()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	<-ctx.Done()
	log.Println("timeout of 5 seconds.")

	log.Println("Server exiting")
}

func (s *server) configureRoutes() {
	rootRoute := s.router.Group("/gowschat")
	rootRoute.GET("/", handleHome)
	rootRoute.POST("/login", func(c *gin.Context) {
		handleLogin(c, s)
	})
	restricted := rootRoute.Group("/chat")
	restricted.GET("/", handleChat)
	restricted.GET("/ws", func(ctx *gin.Context) {
		serveWs(ctx, s)
	})
}

func handleHome(c *gin.Context) {
	view.Login().Render(c.Request.Context(), c.Writer)
}

func handleChat(c *gin.Context) {
	otp, err := c.Cookie(OTP_KEY)
	if err != nil {
		view.Error(err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}
	log.Println("handleChat ::", otp)
	component := view.Chat(otp)
	component.Render(c.Request.Context(), c.Writer)
}
