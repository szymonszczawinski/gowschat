package user

type (
	UserCredentials struct {
		email    string
		password string
	}
	ChatUser struct {
		UserCredentials
	}
)

func NewUserCredentials(email, password string) UserCredentials {
	return UserCredentials{
		email:    email,
		password: password,
	}
}

func NewChatUser(c UserCredentials) *ChatUser {
	return &ChatUser{
		UserCredentials: c,
	}
}
