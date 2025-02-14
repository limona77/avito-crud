package shop

import (
	"avito-crud/internal/client/db"
	"avito-crud/internal/model"
	"avito-crud/internal/repostiory"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrMerchNotFound     = errors.New("merch not found")
)

type shopRepository struct {
	db db.Client
}

func NewShopRepository(db db.Client) repostiory.IShopRepository {
	return &shopRepository{db: db}
}

func (s *shopRepository) GetMerch(ctx context.Context, item string) (*model.Merch, error) {
	op := "shopRepository.GetMerch"
	query := `SELECT * FROM merch WHERE name = $1`

	q := db.Query{
		Name:     "shopRepository.GetMerch",
		QueryRaw: query,
	}

	var merch model.Merch
	err := s.db.DB().QueryRowContext(ctx, q, item).Scan(
		&merch.ID,
		&merch.Name,
		&merch.Price,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.Merch{}, fmt.Errorf("%s: %w", op, ErrMerchNotFound)
		}
		return &model.Merch{}, fmt.Errorf("%s: %w", op, err)
	}
	return &merch, nil
}

func (s *shopRepository) UpdateBalance(ctx context.Context, price, userID int) error {
	op := "shopRepository.UpdateBalance"
	queryUpdateBalance := `
		UPDATE employees
		SET balance = balance - $1
		WHERE id = $2 AND balance >= $1
	`
	qUpdateBalance := db.Query{
		Name:     "shopRepository.UpdateUserBalance",
		QueryRaw: queryUpdateBalance,
	}

	rows, err := s.db.DB().ExecContext(ctx, qUpdateBalance, price, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, ErrInsufficientFunds)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	if rows.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrInsufficientFunds)
	}
	return nil
}
func (s *shopRepository) CreatePurchase(ctx context.Context, userID, merchID int) error {
	op := "shopRepository.InsertPurchase"
	queryInsertPurchase := `
		INSERT INTO purchases (employee_id, merch_id, quantity)
		VALUES ($1, $2, 1)
		ON CONFLICT (employee_id, merch_id)
		DO UPDATE SET quantity = purchases.quantity + 1;
	`
	qInsertPurchase := db.Query{
		Name:     "shopRepository.InsertPurchase",
		QueryRaw: queryInsertPurchase,
	}

	res, err := s.db.DB().ExecContext(ctx, qInsertPurchase, userID, merchID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrInsufficientFunds)
	}
	return nil
}
