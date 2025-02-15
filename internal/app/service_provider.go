package app

import (
	"avito-crud/internal/client/db"
	"avito-crud/internal/client/db/pg"
	"avito-crud/internal/client/db/transaction"
	"avito-crud/internal/closer"
	"avito-crud/internal/config"
	"avito-crud/internal/repostiory"
	authRepo "avito-crud/internal/repostiory/auth"
	infoRepo "avito-crud/internal/repostiory/info"
	shopRepo "avito-crud/internal/repostiory/shop"
	transferRepo "avito-crud/internal/repostiory/transfer"
	"avito-crud/internal/service"
	"avito-crud/internal/service/auth"
	info "avito-crud/internal/service/info"
	"avito-crud/internal/service/shop"
	"avito-crud/internal/service/transfer"
	"avito-crud/internal/utils"
	"context"
	"log"
	"log/slog"
)

type serviceProvider struct {
	httpConfig config.HTTPConfig

	log *slog.Logger

	pgConfig  config.PGConfig
	dbClient  db.Client
	txManager db.TxManager

	tokenService utils.ITokenService

	authService    service.IAuthService
	authRepository repostiory.IAuthRepository

	shopService    service.IShopService
	shopRepository repostiory.IShopRepository

	transferService    service.ITransferService
	transferRepository repostiory.ITransferRepository

	infoRepository repostiory.IinfoRepository
	infoService    service.IInfoService

	jwtConfig config.JWTConfig
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}
func (s *serviceProvider) JWTConfig() config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			log.Fatalf("failed to get jwt config: %s", err.Error())
		}

		s.jwtConfig = cfg
	}

	return s.jwtConfig
}
func (s *serviceProvider) LoggerConfig() *slog.Logger {
	if s.log == nil {
		cfg := config.NewLoggerConfig()

		s.log = cfg
	}

	return s.log
}
func (s *serviceProvider) TokenService() utils.ITokenService {
	if s.tokenService == nil {
		s.tokenService = utils.NewTokenService()
	}
	return s.tokenService
}
func (s *serviceProvider) AuthService(ctx context.Context) service.IAuthService {
	if s.authService == nil {
		s.authService = auth.NewAuthService(s.LoggerConfig(), s.JWTConfig().TTL(), s.AuthRepository(ctx), []byte(s.jwtConfig.Secret()), s.TokenService())
	}

	return s.authService
}
func (s *serviceProvider) AuthRepository(ctx context.Context) repostiory.IAuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepo.NewAuthRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) ShopService(ctx context.Context) service.IShopService {
	if s.shopService == nil {
		s.shopService = shop.NewShopService(s.LoggerConfig(), s.ShopRepository(ctx), s.AuthRepository(ctx), []byte(s.JWTConfig().Secret()), s.TxManager(ctx), s.TokenService())
	}
	return s.shopService
}

func (s *serviceProvider) ShopRepository(ctx context.Context) repostiory.IShopRepository {
	if s.shopRepository == nil {
		s.shopRepository = shopRepo.NewShopRepository(s.DBClient(ctx))
	}
	return s.shopRepository
}

func (s *serviceProvider) TransferService(ctx context.Context) service.ITransferService {
	if s.transferService == nil {
		s.transferService = transfer.NewTransferService(s.LoggerConfig(), s.TransferRepository(ctx), []byte(s.JWTConfig().Secret()), s.TxManager(ctx), s.TokenService())
	}
	return s.transferService
}

func (s *serviceProvider) TransferRepository(ctx context.Context) repostiory.ITransferRepository {
	if s.transferRepository == nil {
		s.transferRepository = transferRepo.NewTransferRepository(s.DBClient(ctx))
	}
	return s.transferRepository
}

func (s *serviceProvider) InfoRepository(ctx context.Context) repostiory.IinfoRepository {
	if s.infoRepository == nil {
		s.infoRepository = infoRepo.NewInfoRepository(s.DBClient(ctx))
	}
	return s.infoRepository
}

func (s *serviceProvider) InfoService(ctx context.Context) service.IInfoService {
	if s.infoService == nil {
		s.infoService = info.NewInfoService(s.LoggerConfig(), s.InfoRepository(ctx), []byte(s.JWTConfig().Secret()), s.TokenService())
	}
	return s.infoService
}
