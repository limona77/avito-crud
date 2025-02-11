package auth

import (
	"avito-crud/internal/repostiory"
	"avito-crud/internal/service"
	"log/slog"
	"time"
)

type auth struct {
	log            *slog.Logger
	tokenTTL       time.Duration
	authRepository repostiory.IAuthRepository
	jwtSecret      []byte
}

func NewAuthService(log *slog.Logger, tokenTTL time.Duration, authRepository repostiory.IAuthRepository, jwtSecret []byte) service.IAuthService {
	return &auth{log: log, tokenTTL: tokenTTL, authRepository: authRepository, jwtSecret: jwtSecret}
}
