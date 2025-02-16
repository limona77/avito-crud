package auth

import (
	"avito-crud/internal/client/db"
	"avito-crud/internal/model"
	"avito-crud/internal/repostiory"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type authRepository struct {
	db db.Client
}

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	balance         = 1000
)

func NewAuthRepository(db db.Client) repostiory.IAuthRepository {
	return &authRepository{db: db}
}

func (a *authRepository) GetUser(ctx context.Context, username string) (*model.User, error) {
	const op = "authRepository.GetUser"

	query := `SELECT * FROM employees WHERE name = $1`

	q := db.Query{
		Name:     "authRepository.GetUser",
		QueryRaw: query,
	}

	var user model.User

	err := a.db.DB().QueryRowContext(ctx, q, username).
		Scan(&user.ID, &user.UserName, &user.Balance, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return &model.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (a *authRepository) CreateUser(ctx context.Context, username, password string) (int, error) {
	const op = "authRepository.CreateUser"

	query := `INSERT INTO employees (name, password, balance) VALUES ($1, $2, $3) RETURNING id`

	q := db.Query{
		Name:     "authRepository.CreateUser",
		QueryRaw: query,
	}

	var id int
	err := a.db.DB().QueryRowContext(ctx, q, username, password, balance).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 - ошибка уникального ограничения в PostgreSQL
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
