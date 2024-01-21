package server

import (
	"errors"
	"gowschat/server/view"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	OTP_KEY          string = "otp"
	USER_SESSION_KEY string = "user_session"
)

var (
	ErrorInvalidSessionToken = errors.New("invalid session token")
	ErrorFailedSaveSession   = errors.New("faield to save session")
)

func handleLoginSubmit(c *gin.Context, s *server) {
	log.Println("login submit")
	password := c.PostForm("username")
	username := c.PostForm("password")

	log.Println("login user", username, "pass", password)
	otp, err := s.authenticator.Login(username, password)
	if err != nil {
		log.Println("login error", err)
		view.LoginError(err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}
	userSessionToken := username + "." + uuid.NewString()
	session := sessions.Default(c)
	session.Set(USER_SESSION_KEY, userSessionToken)
	session.Set(OTP_KEY, otp.Key)
	session.Save()
	c.Writer.Header().Add("HX-Redirect", "/gowschat/chat")
}

func handleLogout(c *gin.Context, s *server) {
	session := sessions.Default(c)
	otp := session.Get(OTP_KEY)
	userSessionId := session.Get(USER_SESSION_KEY)
	if otp == nil {
		c.HTML(http.StatusBadRequest, "error", gin.H{"error": ErrorInvalidSessionToken})
		return
	}
	session.Delete(OTP_KEY)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{"error": ErrorFailedSaveSession})
		return
	}

	userToken, _ := userSessionId.(string)
	userNameAndToken := strings.Split(userToken, ".")
	s.authenticator.Logout(userNameAndToken[0])
	c.Redirect(http.StatusFound, "/gowschat")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func handleRegister(c *gin.Context) {
}

func sessionAuth(c *gin.Context) {
	session := sessions.Default(c)

	if session != nil {
		if session.Get(USER_SESSION_KEY) == nil {
			log.Println("ERROR :: no user session key")
			c.Redirect(http.StatusFound, "/gowschat/login")
			c.Abort()
			return
		}
	}

	c.Next()
}
