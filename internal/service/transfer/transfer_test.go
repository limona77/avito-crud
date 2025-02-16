package transfer

import (
	"avito-crud/internal/model"
	"avito-crud/internal/repostiory/auth"
	"avito-crud/internal/service/shop"
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransfer_SendCoin(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	jwtSecret := []byte("supersecret")
	ctx := context.Background()

	tests := []struct {
		name       string
		token      string
		receiver   string
		amount     int
		setupMocks func(
			tokenService *MockTokenService,
			transferRepo *MockTransferRepository,
			txManager *MockTxManager,
		)
		expectError   bool
		errorContains string
	}{
		{
			name:     "Success",
			token:    "valid-token",
			receiver: "receiver",
			amount:   100,
			setupMocks: func(
				tokenService *MockTokenService,
				transferRepo *MockTransferRepository,
				txManager *MockTxManager,
			) {
				userClaims := &model.UserClaims{
					ID:       1,
					Username: "sender",
				}
				tokenService.
					On("VerifyToken", "valid-token", jwtSecret).
					Return(userClaims, nil)

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
				transferRepo.
					On("CreateTransaction", mock.Anything, "sender", "receiver", 100).
					Return(1, nil)
				// Затем вызывается метод Transfer.
				transferRepo.
					On("Transfer", mock.Anything, "sender", "receiver", 100).
					Return(nil)
			},
			expectError: false,
		},
		{
			name:     "Invalid token",
			token:    "invalid-token",
			receiver: "receiver",
			amount:   100,
			setupMocks: func(
				tokenService *MockTokenService,
				transferRepo *MockTransferRepository,
				txManager *MockTxManager,
			) {
				tokenService.
					On("VerifyToken", "invalid-token", jwtSecret).
					Return(nil, shop.ErrInvalidToken)
			},
			expectError:   true,
			errorContains: shop.ErrInvalidToken.Error(),
		},
		{
			name:     "Repository error on Transfer user not found",
			token:    "valid-token",
			receiver: "receiver",
			amount:   100,
			setupMocks: func(
				tokenService *MockTokenService,
				transferRepo *MockTransferRepository,
				txManager *MockTxManager,
			) {
				userClaims := &model.UserClaims{
					ID:       1,
					Username: "sender",
				}
				tokenService.
					On("VerifyToken", "valid-token", jwtSecret).
					Return(userClaims, nil)

				txManager.
					On(
						"ReadCommitted",
						mock.Anything,
						mock.Anything,
						mock.AnythingOfType("func(context.Context) error"),
					).
					Return(auth.ErrUserNotFound).
					Run(func(args mock.Arguments) {
						callback := args.Get(2).(func(context.Context) error)
						_ = callback(context.Background())
					})

				transferRepo.
					On("CreateTransaction", mock.Anything, "sender", "receiver", 100).
					Return(1, nil)

				// Здесь метод Transfer возвращает ошибку.
				transferRepo.
					On("Transfer", mock.Anything, "sender", "receiver", 100).
					Return(auth.ErrUserNotFound)
			},
			expectError:   true,
			errorContains: auth.ErrUserNotFound.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenServiceMock := new(MockTokenService)
			transferRepoMock := new(MockTransferRepository)
			txManagerMock := new(MockTxManager)

			tt.setupMocks(tokenServiceMock, transferRepoMock, txManagerMock)

			service := NewTransferService(logger, transferRepoMock, jwtSecret, txManagerMock, tokenServiceMock)

			err := service.SendCoin(ctx, tt.token, tt.receiver, tt.amount)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			tokenServiceMock.AssertExpectations(t)
			transferRepoMock.AssertExpectations(t)
			txManagerMock.AssertExpectations(t)
		})
	}
}
