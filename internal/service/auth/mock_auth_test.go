package auth

import (
	"avito-crud/internal/model"
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateToken(user model.User, secretKey []byte, duration time.Duration) (string, error) {
	args := m.Called(user, secretKey, duration)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
	args := m.Called(tokenStr, secretKey)
	// Предполагаем, что если args.Get(0) не nil, то это *model.UserClaims
	if claims := args.Get(0); claims != nil {
		return claims.(*model.UserClaims), args.Error(1)
	}
	return nil, args.Error(1)
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
