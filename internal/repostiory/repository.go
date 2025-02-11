package repostiory

import (
	"avito-crud/internal/model"
	"context"
)

type IAuthRepository interface {
	GetUser(ctx context.Context, username string) (*model.User, error)
	CreateUser(ctx context.Context, username, password string) (int, error)
}
