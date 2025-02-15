package service

import (
	"avito-crud/internal/model"
	"context"
)

type IAuthService interface {
	Auth(ctx context.Context, username, password string) (string, error)
}

type IShopService interface {
	BuyItem(ctx context.Context, token, item string) error
}

type ITransferService interface {
	SendCoin(ctx context.Context, token, receiver string, amount int) error
}

type IInfoService interface {
	GetInfo(ctx context.Context, token string) (*model.UserInfo, error)
}
