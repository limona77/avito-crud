package service

type IAuthService interface {
	Login(username, password string) (string, error)
}
