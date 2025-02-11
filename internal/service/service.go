package service

import "context"

type IAuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
}
