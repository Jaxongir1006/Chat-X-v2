package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/config"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/logger"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres"
	adminRepo "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/admin"
	authRepo "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/auth"
	sessionInfra "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/session"
	redisInfra "github.com/Jaxongir1006/Chat-X-v2/internal/infra/redis"
	redisStore "github.com/Jaxongir1006/Chat-X-v2/internal/infra/redis/store"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/security"
	"github.com/Jaxongir1006/Chat-X-v2/internal/server"
	"github.com/Jaxongir1006/Chat-X-v2/internal/transport/http/auth"
	"github.com/Jaxongir1006/Chat-X-v2/internal/transport/http/middleware"
	"github.com/Jaxongir1006/Chat-X-v2/internal/usecase/adminUsecase"
	authUsecase "github.com/Jaxongir1006/Chat-X-v2/internal/usecase/auth"
)

func Run(cmd string) {
	if cmd == "http" {
		runHttp()
	}
	if cmd == "superuser" {
		createSuperuser()
	}
}

func runHttp() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logger.New(cfg.AppMode)

	ctx := waitForShutdown()

	// init logger

	// init database
	dbPool, err := postgres.New(cfg.PostgresConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize postgres")
		return
	}
	defer func() {
		if err := dbPool.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close db pool")
		}
	}()


	// init redis
	redisPool := redisInfra.NewRedisClient(cfg.RedisConfig)
	if err := redisPool.InitRedis(); err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize redis")
		return
	}
	defer func() {
		if err := redisPool.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close redis pool")
		}
	}()

	// init minio
	// minioStore, err := minio.New(cfg.MinioConfig)
	// if err != nil {
	// 	logger.Fatal().Err(err).Msg("failed to initialized minio")
	// 	return
	// }

	// if err := minioStore.EnsureBucket(ctx); err != nil {
	// 	logger.Fatal().Err(err).Msg("failed to ensure minio bucket")
	// 	return
	// }

	// init infras
	infraSession := sessionInfra.NewSessionRepo(dbPool.DB, logger)

	// init middlewares
	authMiddleware := middleware.NewAuthMiddleware(infraSession, false)

	// init repos
	authRepo := authRepo.NewAuthRepo(dbPool.DB, logger)

	// init services
	hasher := security.NewBcryptHasher(10)
	codeHasher := security.NewHMACHasher("secret")
	redis := redisStore.NewOTPRedisStore(redisPool.Client)
	tokenSrv := security.NewToken(cfg.TokenConfig)

	// init usecases
	authUsecase := authUsecase.NewAuthUsecase(authRepo, infraSession, redis, tokenSrv, hasher, logger, codeHasher)

	// init handlers
	authHandler := auth.NewAuthHandler(authUsecase, logger)

	// init server
	srv := server.NewServer(cfg.Server, authMiddleware, authHandler, logger)

	// start server async
	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to run server")
		}
	}()

	// wait for signal
	<-ctx.Done()
	logger.Info().Msg("Shutting down server...")

	// shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Err(err).Msg("Failed to shutdown server")
	}

	logger.Info().Msg("Shutdown complete")
}

func createSuperuser() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logger.New(cfg.AppMode)

	dbPool, err := postgres.New(cfg.PostgresConfig)
	if err != nil {
		log.Fatalf("Failed to initialize postgres: %v", err)
	}

	defer func() {
		if err := dbPool.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close db pool")
		}
	}()

	repo := adminRepo.NewAdminRepo(dbPool.DB, logger)
	hasher := security.NewBcryptHasher(8)

	usecase := adminUsecase.NewAdminUsecase(repo, hasher)

	err = usecase.CreateSuperuser()
	if err != nil {
		log.Fatalf("Failed to create superuser: %v", err)
	}
}
