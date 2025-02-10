package app

import (
	"avito-crud/internal/client/db"
	"avito-crud/internal/client/db/pg"
	"avito-crud/internal/closer"
	"avito-crud/internal/config"
	"avito-crud/internal/repostiory"
	authRepo "avito-crud/internal/repostiory/auth"
	"avito-crud/internal/service"
	"avito-crud/internal/service/auth"
	"context"
	"log"
	"log/slog"
	"time"
)

type serviceProvider struct {
	httpConfig config.HTTPConfig
	log        *slog.Logger
	pgConfig   config.PGConfig
	dbClient   db.Client

	authService    service.IAuthService
	authRepository repostiory.IAuthRepository

	tokenTTL time.Duration
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

func (s *serviceProvider) AuthService() service.IAuthService {
	if s.authService == nil {
		s.authService = auth.NewAuthService(s.log, s.tokenTTL)
	}

	return s.authService
}
func (s *serviceProvider) NoteRepository(ctx context.Context) repostiory.IAuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepo.NewAuthRepository(s.DBClient(ctx))
	}

	return s.authRepository
}
