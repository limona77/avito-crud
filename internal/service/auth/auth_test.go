package auth

import (
	"avito-crud/internal/model"
	authRepo "avito-crud/internal/repostiory/auth"
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Auth_TableDriven(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	tokenTTL := time.Hour
	jwtSecret := []byte("supersecret")

	tests := []struct {
		name          string
		username      string
		password      string
		setupMocks    func(mockRepo *MockAuthRepository, mockTokenService *MockTokenService)
		expectedToken string
		expectError   bool
	}{
		{
			name:     "Existing user success",
			username: "testuser",
			password: "password",
			setupMocks: func(mockRepo *MockAuthRepository, mockTokenService *MockTokenService) {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
				assert.NoError(t, err)
				existingUser := &model.User{
					ID:       1,
					UserName: "testuser",
					Password: string(hashedPassword),
					Balance:  1000,
				}
				mockRepo.On("GetUser", mock.Anything, "testuser").Return(existingUser, nil)
				mockTokenService.On("GenerateToken", *existingUser, jwtSecret, tokenTTL).
					Return("fake-token", nil)
			},
			expectedToken: "fake-token",
			expectError:   false,
		},
		{
			name:     "New user success",
			username: "newuser",
			password: "password",
			setupMocks: func(mockRepo *MockAuthRepository, mockTokenService *MockTokenService) {
				mockRepo.
					On("GetUser", mock.Anything, "newuser").
					Return(&model.User{}, authRepo.ErrUserNotFound)
				mockRepo.
					On("CreateUser", mock.Anything, "newuser", mock.AnythingOfType("string")).
					Return(2, nil)
				newUser := &model.User{
					ID:       2,
					UserName: "newuser",
					Balance:  0,
				}
				mockTokenService.On("GenerateToken", *newUser, jwtSecret, tokenTTL).
					Return("new-fake-token", nil)
			},
			expectedToken: "new-fake-token",
			expectError:   false,
		},
		{
			name:     "Invalid credentials",
			username: "testuser",
			password: "wrongpassword",
			setupMocks: func(mockRepo *MockAuthRepository, mockTokenService *MockTokenService) {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
				assert.NoError(t, err)
				existingUser := &model.User{
					ID:       1,
					UserName: "testuser",
					Password: string(hashedPassword),
					Balance:  1000,
				}
				mockRepo.On("GetUser", mock.Anything, "testuser").Return(existingUser, nil)
			},
			expectedToken: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAuthRepository)
			mockTokenService := new(MockTokenService)

			tt.setupMocks(mockRepo, mockTokenService)

			authService := NewAuthService(logger, tokenTTL, mockRepo, jwtSecret, mockTokenService)

			token, err := authService.Auth(context.Background(), tt.username, tt.password)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, token)
			}
			mockRepo.AssertExpectations(t)
			mockTokenService.AssertExpectations(t)
		})
	}
}
