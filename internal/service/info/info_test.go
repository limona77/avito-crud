package info

import (
	"avito-crud/internal/model"
	"avito-crud/internal/service/shop"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInfoService_GetInfo(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	jwtSecret := []byte("supersecret")
	ctx := context.Background()

	tests := []struct {
		name          string
		token         string
		setupMocks    func(tokenService *MockTokenService, repo *MockInfoRepositoryRepository)
		expectedInfo  *model.UserInfo
		expectError   bool
		errorContains string
	}{
		{
			name:  "Valid token and repository success",
			token: "valid-token",
			setupMocks: func(tokenService *MockTokenService, repo *MockInfoRepositoryRepository) {
				userClaims := &model.UserClaims{
					ID:       1,
					Username: "testuser",
				}
				tokenService.
					On("VerifyToken", "valid-token", jwtSecret).
					Return(userClaims, nil)

				repo.
					On("GetBalance", mock.Anything, userClaims.ID).
					Return(1000, nil)

				inventory := []*model.UserInventory{
					{Type: "item1", Amount: 2},
					{Type: "item2", Amount: 1},
				}
				repo.
					On("GetUserInventory", mock.Anything, userClaims.ID).
					Return(inventory, nil)

				received := []*model.Received{
					{Amount: 500},
				}
				repo.
					On("GetReceivedTransactions", mock.Anything, userClaims.Username).
					Return(received, nil)

				sent := []*model.Sent{
					{Amount: 200},
				}
				repo.
					On("GetSentTransactions", mock.Anything, userClaims.Username).
					Return(sent, nil)
			},
			expectedInfo: &model.UserInfo{
				Coins: 1000,
				Inventory: []*model.UserInventory{
					{Type: "item1", Amount: 2},
					{Type: "item2", Amount: 1},
				},
				CoinHistory: model.CoinHistory{
					Received: []*model.Received{
						{Amount: 500},
					},
					Sent: []*model.Sent{
						{Amount: 200},
					},
				},
			},
			expectError: false,
		},
		{
			name:  "Invalid token",
			token: "invalid-token",
			setupMocks: func(tokenService *MockTokenService, repo *MockInfoRepositoryRepository) {
				tokenService.
					On("VerifyToken", "invalid-token", jwtSecret).
					Return(nil, shop.ErrInvalidToken)
			},
			expectedInfo:  &model.UserInfo{},
			expectError:   true,
			errorContains: shop.ErrInvalidToken.Error(),
		},
		{
			name:  "Repository error on GetBalance",
			token: "valid-token",
			setupMocks: func(tokenService *MockTokenService, repo *MockInfoRepositoryRepository) {
				userClaims := &model.UserClaims{
					ID:       1,
					Username: "testuser",
				}
				tokenService.
					On("VerifyToken", "valid-token", jwtSecret).
					Return(userClaims, nil)
				repo.
					On("GetBalance", mock.Anything, userClaims.ID).
					Return(0, errors.New("db error"))
			},
			expectedInfo:  nil,
			expectError:   true,
			errorContains: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenServiceMock := new(MockTokenService)
			repoMock := new(MockInfoRepositoryRepository)

			tt.setupMocks(tokenServiceMock, repoMock)

			infoService := NewInfoService(logger, repoMock, jwtSecret, tokenServiceMock)

			userInfo, err := infoService.GetInfo(ctx, tt.token)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedInfo, userInfo)
			}

			repoMock.AssertExpectations(t)
			tokenServiceMock.AssertExpectations(t)
		})
	}
}
