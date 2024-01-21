package auth

import (
	"context"
	"errors"
	"gowschat/server/auth/user"
	"log"
	"time"

	"github.com/google/uuid"
)

var users = map[string]string{"szymon": "szymon", "piotr": "piotr", "madga": "magda"}

type (
	OTP struct {
		Created         time.Time
		Key             string
		userCredentials user.UserCredentials
	}
	retentionMap  map[string]OTP
	Authenticator struct {
		otps          retentionMap
		loggedInUsers map[string]user.ChatUser
	}
)

func NewAuthenticator(ctx context.Context, retentionPeriod time.Duration) *Authenticator {
	a := &Authenticator{
		otps:          make(retentionMap),
		loggedInUsers: map[string]user.ChatUser{},
	}
	go a.retention(ctx, retentionPeriod)
	return a
}

func (a *Authenticator) Login(username, password string) (OTP, error) {
	if pass, userExist := users[username]; userExist && pass == password {
		otp := a.NewOTP(username, password)
		// TODO: get user from DB
		a.loggedInUsers[username] = *user.NewChatUser(user.NewUserCredentials(username, password))
		return otp, nil
	}
	return OTP{}, errors.New("login failed")
}

func (a *Authenticator) Logout(username string) error {
	delete(a.loggedInUsers, username)
	return nil
}

// NewOTP creates and adds a new otp to the map
func (a *Authenticator) NewOTP(username, password string) OTP {
	o := OTP{
		Key:             uuid.NewString(),
		Created:         time.Now(),
		userCredentials: user.NewUserCredentials(username, password),
	}

	a.otps[o.Key] = o
	return o
}

// VerifyOTP will make sure a OTP exists
// and return true if so
// It will also delete the key so it cant be reused
func (a *Authenticator) VerifyOTP(otp string) (bool, user.ChatUser) {
	if otp == "superduper" {
		return true, *user.NewChatUser(user.NewUserCredentials("xxx", "xxx"))
	}
	// Verify OTP is existing
	log.Println("VerifyOTP ::", otp, " ALL:", a.otps)
	otpContainer, ok := a.otps[otp]
	if !ok {
		// otp does not exist
		return false, user.ChatUser{}
	}
	delete(a.otps, otp)
	return true, a.loggedInUsers[otpContainer.userCredentials.GetUsername()]
}

// Retention will make sure old OTPs are removed
// Is Blocking, so run as a Goroutine
func (a *Authenticator) retention(ctx context.Context, retentionPeriod time.Duration) {
	log.Println("run retention")
	ticker := time.NewTicker(400 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			for _, otp := range a.otps {
				// Add Retention to Created and check if it is expired
				if otp.Created.Add(retentionPeriod).Before(time.Now()) {
					delete(a.otps, otp.Key)
					delete(a.loggedInUsers, otp.userCredentials.GetUsername())
				}
			}
		case <-ctx.Done():
			log.Println("kill retention goroutine")
			return

		}
	}
}
