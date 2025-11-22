package dto

type Auth interface {
	GetUser() string
	GetPassword() string
}

type AuthDTO struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func (a *AuthDTO) GetUser() string {
	return a.User
}

func (a *AuthDTO) GetPassword() string {
	return a.Password
}
