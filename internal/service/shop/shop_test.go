package shop

import (
	"avito-crud/internal/model"
	shopRepo "avito-crud/internal/repostiory/shop"
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShop_BuyItem(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	jwtSecret := []byte("supersecret")
	ctx := context.Background()

	tests := []struct {
		name       string
		token      string
		item       string
		setupMocks func(
			tokenService *MockTokenService,
			shopRepo *MockShopRepository,
			authRepoMock *MockAuthRepository,
			txManager *MockTxManager,
		)
		expectError   bool
		errorContains string
	}{
		{
			name:  "Success",
			item:  "item1",
			token: "valid-token",
			setupMocks: func(
				tokenService *MockTokenService,
				shopRepo *MockShopRepository,
				authRepoMock *MockAuthRepository,
				txManager *MockTxManager,
			) {
				txManager.
					On(
						"ReadCommitted",
						mock.Anything,
						mock.Anything,
						mock.AnythingOfType("func(context.Context) error"),
					).
					Return(nil).
					Run(func(args mock.Arguments) {
						callback := args.Get(2).(func(context.Context) error)
						_ = callback(context.Background())
					})
				authRepoMock.
					On("GetUser", mock.Anything, "testuser").
					Return(&model.User{ID: 1, UserName: "testuser", Balance: 1000}, nil)
				userClaims := &model.UserClaims{
					ID:       1,
					Username: "testuser",
				}
				tokenService.
					On("VerifyToken", "valid-token", jwtSecret).
					Return(userClaims, nil)

				shopRepo.
					On("GetMerch", mock.Anything, "item1").
					Return(&model.Merch{ID: 1, Name: "item1", Price: 100}, nil)
				shopRepo.
					On("UpdateBalance", mock.Anything, 100, 1).
					Return(nil)
				shopRepo.
					On("CreatePurchase", mock.Anything, 1, 1).
					Return(nil)
			},
			expectError: false,
		},
		{
			name:  "Invalid token",
			item:  "item1",
			token: "invalid-token",
			setupMocks: func(
				tokenService *MockTokenService,
				shopRepo *MockShopRepository,
				authRepoMock *MockAuthRepository,
				txManager *MockTxManager,
			) {
				tokenService.
					On("VerifyToken", "invalid-token", jwtSecret).
					Return(nil, ErrInvalidToken)
			},
			expectError:   true,
			errorContains: ErrInvalidToken.Error(),
		},
		{
			name:  "Repository error on UpdateBalance",
			item:  "item1",
			token: "valid-token",
			setupMocks: func(
				tokenService *MockTokenService,
				shopRepoMock *MockShopRepository,
				authRepoMock *MockAuthRepository,
				txManager *MockTxManager,
			) {
				txManager.
					On(
						"ReadCommitted",
						mock.Anything,
						mock.Anything,
						mock.AnythingOfType("func(context.Context) error"),
					).
					Return(shopRepo.ErrInsufficientFunds).
					Run(func(args mock.Arguments) {
						callback := args.Get(2).(func(context.Context) error)
						_ = callback(context.Background())
					})
				authRepoMock.
					On("GetUser", mock.Anything, "testuser").
					Return(&model.User{ID: 1, UserName: "testuser", Balance: 1000}, nil)
				userClaims := &model.UserClaims{
					ID:       1,
					Username: "testuser",
				}
				tokenService.
					On("VerifyToken", "valid-token", jwtSecret).
					Return(userClaims, nil)
				shopRepoMock.
					On("GetMerch", mock.Anything, "item1").
					Return(&model.Merch{ID: 1, Name: "item1", Price: 100}, nil)
				shopRepoMock.
					On("UpdateBalance", mock.Anything, 100, 1).
					Return(shopRepo.ErrInsufficientFunds)
			},
			expectError:   true,
			errorContains: shopRepo.ErrInsufficientFunds.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenServiceMock := new(MockTokenService)
			shopRepoMock := new(MockShopRepository)
			authRepoMock := new(MockAuthRepository)
			txManagerMock := new(MockTxManager)
			tt.setupMocks(tokenServiceMock, shopRepoMock, authRepoMock, txManagerMock)

			shopService := NewShopService(
				logger,
				shopRepoMock,
				authRepoMock,
				jwtSecret,
				txManagerMock,
				tokenServiceMock,
			)

			err := shopService.BuyItem(ctx, tt.token, tt.item)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
			txManagerMock.AssertExpectations(t)
			authRepoMock.AssertExpectations(t)
			shopRepoMock.AssertExpectations(t)
			tokenServiceMock.AssertExpectations(t)
		})
	}
}
