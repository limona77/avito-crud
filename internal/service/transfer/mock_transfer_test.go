package transfer

import (
	"avito-crud/internal/model"
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
)

type MockTransferRepository struct {
	mock.Mock
}

func (t *MockTransferRepository) CreateTransaction(
	ctx context.Context,
	sender, receiver string,
	amount int,
) (int, error) {
	args := t.Called(ctx, sender, receiver, amount)
	return args.Int(0), args.Error(1)
}

func (t *MockTransferRepository) Transfer(
	ctx context.Context,
	sender, receiver string,
	amount int,
) error {
	args := t.Called(ctx, sender, receiver, amount)
	return args.Error(0)
}

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateToken(
	user model.User,
	secretKey []byte,
	duration time.Duration,
) (string, error) {
	args := m.Called(user, secretKey, duration)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) VerifyToken(
	tokenStr string,
	secretKey []byte,
) (*model.UserClaims, error) {
	args := m.Called(tokenStr, secretKey)
	if claims := args.Get(0); claims != nil {
		return claims.(*model.UserClaims), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockTxManager struct {
	mock.Mock
}

func (m *MockTxManager) ReadCommitted(ctx context.Context, f func(ctx context.Context) error) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.Called(ctx, txOpts, f).Error(0)
}
