package dto

type Auth interface {
	GetUser() string
	GetPassword() string
}

type authImpl struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func (a *authImpl) GetUser() string {
	return a.User
}

func (a *authImpl) GetPassword() string {
	return a.Password
}

func NewAuth(user, password string) Auth {
	return &authImpl{User: user, Password: password}
}
