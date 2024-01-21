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

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const (
	// TODO: replace with Env variable
	MySecret = "mysecret"
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
		router:        createGinRouter(),
	}
}

func createGinRouter() *gin.Engine {
	engine := gin.Default()
	// cookieOptions := sessions.Options{
	// 	Path:     "/",
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteLaxMode,
	// 	Domain:   "gowschat",
	// 	MaxAge:   60 * 5,
	// }
	cookieStore := cookie.NewStore([]byte(os.Getenv(MySecret)))
	// cookieStore.Options(cookieOptions)
	engine.Use(sessions.Sessions("mysession", cookieStore))
	return engine
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
	rootRoute.GET("/login", handleLogin)
	rootRoute.POST("/login", func(c *gin.Context) {
		handleLoginSubmit(c, s)
	})
	restricted := rootRoute.Group("/chat")
	restricted.Use(sessionAuth)
	restricted.GET("/", handleChat)
	restricted.GET("/ws", func(ctx *gin.Context) {
		serveWs(ctx, s)
	})
}

func handleHome(c *gin.Context) {
	view.Home().Render(c.Request.Context(), c.Writer)
}

func handleLogin(c *gin.Context) {
	view.Login().Render(c.Request.Context(), c.Writer)
}

func handleChat(c *gin.Context) {
	session := sessions.Default(c)
	otp := session.Get(OTP_KEY)
	if otp == nil {
		log.Println("ERROR: no OTP")
		view.Error("No OTP found").Render(c.Request.Context(), c.Writer)
		return
	}
	otpString, _ := otp.(string)
	component := view.Chat(otpString)
	component.Render(c.Request.Context(), c.Writer)
}
