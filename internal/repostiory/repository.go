package repostiory

import (
	"avito-crud/internal/model"
	"context"
)

type IAuthRepository interface {
	GetUser(ctx context.Context, username string) (*model.User, error)
	CreateUser(ctx context.Context, username, password string) (int, error)
}

type IShopRepository interface {
	GetMerch(ctx context.Context, name string) (*model.Merch, error)
	UpdateBalance(ctx context.Context, price, userID int) error
	CreatePurchase(ctx context.Context, userID, merchID int) error
}
