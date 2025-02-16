package jwtToken

import (
	"avito-crud/internal/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// TokenService описывает интерфейс для генерации и проверки токенов.
type ITokenService interface {
	GenerateToken(user model.User, secretKey []byte, duration time.Duration) (string, error)
	VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error)
}

// DefaultTokenService — реализация интерфейса TokenService по умолчанию.
type TokenSevice struct{}

func NewTokenService() ITokenService {
	return &TokenSevice{}
}

// GenerateToken генерирует JWT-токен на основе информации о пользователе.
func (d *TokenSevice) GenerateToken(info model.User, secretKey []byte, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},

		Username: info.UserName,
		Balance:  info.Balance,
		ID:       info.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// VerifyToken проверяет валидность токена и возвращает данные из claims.
func (d *TokenSevice) VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Проверяем, что используется ожидаемый метод подписи.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.Errorf("unexpected token signing method")
			}
			return secretKey, nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.Errorf("invalid token claims")
	}

	return claims, nil
}
