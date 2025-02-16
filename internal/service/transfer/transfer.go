package transfer

import (
	"avito-crud/internal/client/db"
	"avito-crud/internal/repostiory"
	"avito-crud/internal/repostiory/auth"
	shopRepo "avito-crud/internal/repostiory/shop"
	"avito-crud/internal/service"
	"avito-crud/internal/service/shop"
	"avito-crud/internal/utils/jwtToken"
	"avito-crud/pkg/logger/sl"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var ErrSameUser = errors.New("sender and receiver are the same")

type transfer struct {
	log                *slog.Logger
	transferRepository repostiory.ITransferRepository
	txManager          db.TxManager
	tokenService       jwtToken.ITokenService
	jwtSecret          []byte
}

func NewTransferService(
	log *slog.Logger,
	transactionRepository repostiory.ITransferRepository,
	jwtSecret []byte,
	txManager db.TxManager,
	tokenService jwtToken.ITokenService,
) service.ITransferService {
	return &transfer{
		log:                log,
		transferRepository: transactionRepository,
		jwtSecret:          jwtSecret,
		txManager:          txManager,
		tokenService:       tokenService,
	}
}

func (t *transfer) SendCoin(ctx context.Context, token, receiver string, amount int) error {
	const op = "transfer.Transfer"
	log := t.log.With(
		slog.String("op", op),
		slog.String("sender, receiver", receiver),
	)
	log.Info("verifying token")
	userClaims, err := t.tokenService.VerifyToken(token, t.jwtSecret)
	if err != nil {
		t.log.Warn("failed to verify token", sl.Err(err))
		return fmt.Errorf("%s: %w", op, shop.ErrInvalidToken)
	}
	if userClaims.Username == receiver {
		t.log.Warn("sender and receiver are the same")
		return fmt.Errorf("%s: %w", op, ErrSameUser)
	}
	log.Info("starting transaction")
	err = t.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		log.Info("attempting to transfer")
		var errTx error
		_, errTx = t.transferRepository.CreateTransaction(ctx, userClaims.Username, receiver, amount)
		if errTx != nil {
			return fmt.Errorf("%s: %w", op, errTx)
		}

		log.Info("attempting to update balance")
		errTx = t.transferRepository.Transfer(ctx, userClaims.Username, receiver, amount)
		if errTx != nil {
			if errors.Is(errTx, shopRepo.ErrInsufficientFunds) {
				t.log.Warn("insufficient funds", sl.Err(errTx))
				return fmt.Errorf("%s: %w", op, shopRepo.ErrInsufficientFunds)
			}
			if errors.Is(errTx, auth.ErrUserNotFound) {
				t.log.Warn("user not found", sl.Err(errTx))
				return fmt.Errorf("%s: %w", op, auth.ErrUserNotFound)
			}
			return fmt.Errorf("%s: %w", op, errTx)
		}
		return nil
	})
	if err != nil {
		t.log.Warn("failed to send Coin", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("end transaction")

	return nil
}
