package server

import (
	"context"
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
	ctx    context.Context
	chat   *chat.ChatServer
	router *gin.Engine
}

func newServer(ctx context.Context) *server {
	return &server{
		ctx:    ctx,
		chat:   chat.NewChatServer(ctx),
		router: gin.Default(),
	}
}

func (s *server) run() {
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
	s.router.POST("/login", func(c *gin.Context) {
		handleLogin(c, s)
	})
	s.router.GET("/", handleHome)
	restricted := s.router.Group("/chat")
	restricted.GET("/", handleChat)
	restricted.GET("/ws", s.chat.ServeWs)
}

func handleHome(c *gin.Context) {
	view.Login().Render(c.Request.Context(), c.Writer)
}

func handleChat(c *gin.Context) {
	otp, err := c.Cookie("otp")
	if err != nil {
		view.Error(err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}
	log.Println("handleChat ::", otp)
	component := view.Chat(otp)
	component.Render(c.Request.Context(), c.Writer)
}

func handleLogin(c *gin.Context, s *server) {
	log.Println("login submit")
	password := c.PostForm("username")
	username := c.PostForm("password")

	log.Println("login user", username, "pass", password)
	otp, err := s.chat.Login(username, password)
	if err != nil {
		log.Println("login error", err)
		view.LoginError(err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}
	log.Println("login otp", otp.Key)
	c.SetCookie("otp", otp.Key, 120, "", c.Request.URL.Hostname(), false, false)
	c.Writer.Header().Add("HX-Redirect", "/chat")
}
