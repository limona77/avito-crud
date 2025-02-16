package info

import (
	"avito-crud/internal/model"
	"avito-crud/internal/repostiory"
	"avito-crud/internal/repostiory/auth"
	"avito-crud/internal/service"
	"avito-crud/internal/service/shop"
	"avito-crud/internal/utils/jwtToken"
	"avito-crud/pkg/logger/sl"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type info struct {
	log            *slog.Logger
	infoRepository repostiory.IinfoRepository
	jwtSecret      []byte
	tokenService   jwtToken.ITokenService
}

func NewInfoService(
	log *slog.Logger,
	infoRepository repostiory.IinfoRepository,
	jwtSecret []byte,
	tokenService jwtToken.ITokenService,
) service.IInfoService {
	return &info{
		log:            log,
		infoRepository: infoRepository,
		jwtSecret:      jwtSecret,
		tokenService:   tokenService,
	}
}

func (i *info) GetInfo(ctx context.Context, token string) (*model.UserInfo, error) {
	const op = "info.GetInfo"

	log := i.log.With(
		slog.String("op", op),
	)
	log.Info("verifying token")
	userClaims, err := i.tokenService.VerifyToken(token, i.jwtSecret)
	if err != nil {
		i.log.Warn("failed to verify token", sl.Err(err))
		return &model.UserInfo{}, fmt.Errorf("%s: %w", op, shop.ErrInvalidToken)
	}

	balance, err := i.infoRepository.GetBalance(ctx, userClaims.ID)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			i.log.Warn("user not found", sl.Err(err))
			return &model.UserInfo{}, fmt.Errorf("%s: %w", op, auth.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	inventory, err := i.infoRepository.GetUserInventory(ctx, userClaims.ID)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			i.log.Warn("user not found", sl.Err(err))
			return &model.UserInfo{}, fmt.Errorf("%s: %w", op, auth.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	received, err := i.infoRepository.GetReceivedTransactions(ctx, userClaims.Username)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			i.log.Warn("user not found", sl.Err(err))
			return &model.UserInfo{}, fmt.Errorf("%s: %w", op, auth.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	sent, err := i.infoRepository.GetSentTransactions(ctx, userClaims.Username)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			i.log.Warn("user not found", sl.Err(err))
			return &model.UserInfo{}, fmt.Errorf("%s: %w", op, auth.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userInfo := &model.UserInfo{
		Coins:     balance,
		Inventory: inventory,
		CoinHistory: model.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}

	return userInfo, nil
}
