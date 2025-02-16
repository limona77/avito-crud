package shop

import (
	"avito-crud/internal/model"
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
)

type MockShopRepository struct {
	mock.Mock
}

func (m *MockShopRepository) GetMerch(ctx context.Context, name string) (*model.Merch, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*model.Merch), args.Error(1)
}

func (m *MockShopRepository) UpdateBalance(ctx context.Context, price, userID int) error {
	args := m.Called(ctx, price, userID)
	return args.Error(0)
}

func (m *MockShopRepository) CreatePurchase(ctx context.Context, userID, merchID int) error {
	args := m.Called(ctx, userID, merchID)
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

func (m *MockTokenService) VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
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

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) GetUser(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, username, password string) (int, error) {
	args := m.Called(ctx, username, password)
	return args.Int(0), args.Error(1)
}
