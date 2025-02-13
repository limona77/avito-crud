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

type ITransferRepository interface {
	CreateTransaction(ctx context.Context, sender, receiver string, amount int) (int, error)
	UpdateBalance(ctx context.Context, sender, receiver string, amount int) error
	CheckBalance(ctx context.Context, sender string, amount int) (bool, error)
}
