package auth

import (
	"avito-crud/internal/service"

	"log/slog"
	"time"
)

type auth struct {
	log      *slog.Logger
	tokenTTL time.Duration
}

func (a auth) Login(username, password string) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", username),
	)

	log.Info("login user")
	return "", nil
}

func NewAuthService(log *slog.Logger, tokenTTL time.Duration) service.IAuthService {
	return &auth{log: log, tokenTTL: tokenTTL}
}
