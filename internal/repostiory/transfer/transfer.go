package transfer

import (
	"avito-crud/internal/client/db"
	"avito-crud/internal/repostiory"
	"avito-crud/internal/repostiory/auth"
	"avito-crud/internal/repostiory/shop"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
)

type transferRepository struct {
	db db.Client
}

func NewTransferRepository(db db.Client) repostiory.ITransferRepository {
	return &transferRepository{db: db}
}

func (t *transferRepository) CreateTransaction(ctx context.Context, sender, receiver string, amount int) (int, error) {
	const op = "transactionRepository.CreateTransaction"
	query := `INSERT INTO transactions (sender, receiver, amount) VALUES ($1, $2, $3) RETURNING id`

	q := db.Query{
		Name:     "transferRepository.CreateTransaction",
		QueryRaw: query,
	}
	var id int
	err := t.db.DB().QueryRowContext(ctx, q, sender, receiver, amount).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return 0, fmt.Errorf("%s: %w", op, auth.ErrUserNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (t *transferRepository) Transfer(ctx context.Context, sender, receiver string, amount int) error {
	const op = "transactionRepository.UpdateBalance"
	query := `WITH deduct AS (
    				UPDATE employees
    				SET balance = balance - $1
    				WHERE name = $2 AND balance >= $1
    				RETURNING id
						)
						UPDATE employees
						SET balance = balance + $1
						WHERE name = $3 AND EXISTS (SELECT 1 FROM deduct);
						`
	q := db.Query{
		Name:     "transferRepository.UpdateBalance",
		QueryRaw: query,
	}

	res, err := t.db.DB().ExecContext(ctx, q, amount, sender, receiver)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 - ошибка уникального ограничения в PostgreSQL
			return fmt.Errorf("%s: %w", op, auth.ErrUserNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, shop.ErrInsufficientFunds)
	}
	return nil
}
