package auth

import (
	repo "avito-crud/internal/repostiory/auth"
	"avito-crud/internal/utils"
	"avito-crud/pkg/logger/sl"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (a *auth) Login(ctx context.Context, username, password string) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", username),
	)

	log.Info("attempting to login user")

	user, err := a.authRepository.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))
			passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				a.log.Error("failed to generate password hash", sl.Err(err))

				return "", fmt.Errorf("%s: %w", op, err)
			}
			_, err = a.authRepository.CreateUser(ctx, username, string(passHash))
			if err != nil {
				a.log.Error("failed to create user", sl.Err(err))
				if errors.Is(err, repo.ErrUserExists) {
					return "", fmt.Errorf("%s: %w", op, repo.ErrUserExists)
				}
				return "", err
			}
			token, err := utils.GenerateToken(*user, a.jwtSecret, a.tokenTTL)
			if err != nil {
				a.log.Error("failed to generate token", sl.Err(err))

				return "", fmt.Errorf("%s: %w", op, err)
			}

			return token, nil
		}
		a.log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)

	} else if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := utils.GenerateToken(*user, a.jwtSecret, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}
