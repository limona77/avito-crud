package shop

import (
	"avito-crud/internal/repostiory"
	authRepo "avito-crud/internal/repostiory/auth"
	shopRepo "avito-crud/internal/repostiory/shop"
	"avito-crud/internal/service"
	"avito-crud/internal/utils"
	"avito-crud/pkg/logger/sl"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type shop struct {
	log            *slog.Logger
	shopRepository repostiory.IShopRepository
	authRepository repostiory.IAuthRepository
	jwtSecret      []byte
}

var (
	ErrInvalidToken = errors.New("invalid token")
)

func NewShopService(log *slog.Logger, shopRepository repostiory.IShopRepository, authRepository repostiory.IAuthRepository, jwtSecret []byte) service.IShopService {
	return &shop{log: log, shopRepository: shopRepository, authRepository: authRepository, jwtSecret: jwtSecret}
}

func (s *shop) BuyItem(ctx context.Context, token, item string) error {
	const op = "shop.BuyItem"
	log := s.log.With(
		slog.String("op", op),
		slog.String("item", item),
	)

	log.Info("verifying token")
	userClaims, err := utils.VerifyToken(token, s.jwtSecret)
	if err != nil {
		s.log.Warn("failed to verify token", sl.Err(err))
		return fmt.Errorf("%s: %w", op, ErrInvalidToken)
	}

	// TODO: Cover with Transaction
	log.Info("attempting to get user")
	user, err := s.authRepository.GetUser(ctx, userClaims.Username)
	if err != nil {
		if errors.Is(err, authRepo.ErrUserNotFound) {
			s.log.Warn("user not found", sl.Err(err))
			return authRepo.ErrUserNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("attempting to get merch")
	merch, err := s.shopRepository.GetMerch(ctx, item)
	if err != nil {
		if errors.Is(err, shopRepo.ErrMerchNotFound) {
			s.log.Warn("merch not found", sl.Err(err))
			return shopRepo.ErrMerchNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if user.Balance < merch.Price {
		return fmt.Errorf("%s: %w", op, shopRepo.ErrInsufficientFunds)
	}
	log.Info("attempting to update balance")
	err = s.shopRepository.UpdateBalance(ctx, merch.Price, user.ID)
	if err != nil {
		if errors.Is(err, shopRepo.ErrInsufficientFunds) {
			s.log.Warn("insufficient funds", sl.Err(err))
			return shopRepo.ErrInsufficientFunds
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	err = s.shopRepository.CreatePurchase(ctx, user.ID, merch.ID)
	if err != nil {
		if errors.Is(err, shopRepo.ErrInsufficientFunds) {
			s.log.Warn("insufficient funds", sl.Err(err))
			return shopRepo.ErrInsufficientFunds
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
