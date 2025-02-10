package auth

import (
	"avito-crud/internal/client/db"
	"avito-crud/internal/repostiory"
)

type authRepository struct {
	db db.Client
}

func NewAuthRepository(db db.Client) repostiory.IAuthRepository {
	return &authRepository{db: db}
}
