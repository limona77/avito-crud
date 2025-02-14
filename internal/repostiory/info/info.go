package info

import (
	"avito-crud/internal/client/db"
	"avito-crud/internal/model"
	"avito-crud/internal/repostiory"
	"avito-crud/internal/repostiory/auth"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type infoRepository struct {
	db db.Client
}

func NewInfoRepository(db db.Client) repostiory.IinfoRepository {
	return &infoRepository{db: db}
}
func (i *infoRepository) GetUserInventory(ctx context.Context, userID int) ([]*model.UserInventory, error) {
	const op = "infoRepository.GetUserInventory"

	query := `SELECT m.name AS type, SUM(p.quantity) AS amount
						FROM purchases p
						JOIN merch m ON p.merch_id = m.id
						WHERE p.employee_id = $1
						GROUP BY m.name;`

	q := db.Query{
		Name:     "infoRepository.GetUserInventory",
		QueryRaw: query,
	}

	var userInventory []*model.UserInventory
	err := i.db.DB().ScanAllContext(ctx, &userInventory, q, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*model.UserInventory{}, fmt.Errorf("%s: %w", op, auth.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return userInventory, nil
}

func (i *infoRepository) GetReceivedTransactions(ctx context.Context, user string) ([]*model.Received, error) {
	const op = "infoRepository.GetReceivedTransactions"

	query := `SELECT transactions.sender , transactions.amount
						FROM transactions
						WHERE transactions.receiver = $1;
	`

	q := db.Query{
		Name:     "infoRepository.GetReceivedTransactions",
		QueryRaw: query,
	}

	var received []*model.Received
	err := i.db.DB().ScanAllContext(ctx, &received, q, user)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*model.Received{}, fmt.Errorf("%s: %w", op, auth.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return received, nil
}

func (i *infoRepository) GetSentTransactions(ctx context.Context, user string) ([]*model.Sent, error) {
	const op = "infoRepository.GetSentTransactions"

	query := `SELECT transactions.receiver, transactions.amount
						FROM transactions
						WHERE transactions.sender = $1;`

	q := db.Query{
		Name:     "infoRepository.GetSentTransactions",
		QueryRaw: query,
	}

	var sent []*model.Sent
	err := i.db.DB().ScanAllContext(ctx, &sent, q, user)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*model.Sent{}, fmt.Errorf("%s: %w", op, auth.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sent, nil
}

func (i *infoRepository) GetBalance(ctx context.Context, userID int) (int, error) {
	const op = "infoRepository.GetBalance"

	query := `SELECT balance
						FROM employees
						WHERE id = $1;`

	q := db.Query{
		Name:     "infoRepository.GetBalance",
		QueryRaw: query,
	}

	var balance int

	err := i.db.DB().ScanOneContext(ctx, &balance, q, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", op, auth.ErrUserNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return balance, nil
}
