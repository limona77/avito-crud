package service

import "context"

type IAuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
}

type IShopService interface {
	BuyItem(ctx context.Context, token, item string) error
}
