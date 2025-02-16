package info

import (
	"avito-crud/internal/model"
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockInfoRepositoryRepository struct {
	mock.Mock
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

func (m *MockInfoRepositoryRepository) GetUserInventory(
	ctx context.Context,
	userID int,
) ([]*model.UserInventory, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.UserInventory), args.Error(1)
}

func (m *MockInfoRepositoryRepository) GetReceivedTransactions(
	ctx context.Context,
	user string,
) ([]*model.Received, error) {
	args := m.Called(ctx, user)
	return args.Get(0).([]*model.Received), args.Error(1)
}

func (m *MockInfoRepositoryRepository) GetSentTransactions(
	ctx context.Context,
	user string,
) ([]*model.Sent, error) {
	args := m.Called(ctx, user)
	return args.Get(0).([]*model.Sent), args.Error(1)
}

func (m *MockInfoRepositoryRepository) GetBalance(ctx context.Context, userID int) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}
