package config

import (
	"errors"
	"os"
	"time"
)

const (
	jwtSecretEnvName = "JWT_SECRET"
	tokenTTl         = "TOKEN_TTL"
)

type JWTConfig interface {
	Secret() string
	TTL() time.Duration
}

type jwtConfig struct {
	secret string
	ttl    time.Duration
}

func NewJWTConfig() (JWTConfig, error) {
	secret := os.Getenv(jwtSecretEnvName)
	ttlStr := os.Getenv(tokenTTl)
	if secret == "" {
		return nil, errors.New("JWT_SECRET not found in environment")
	}
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		ttl = 24 * time.Hour
	}
	return &jwtConfig{secret: secret, ttl: ttl}, nil
}

func (j *jwtConfig) Secret() string {
	return j.secret
}

func (j *jwtConfig) TTL() time.Duration {
	return j.ttl
}
