package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/config"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/minio"
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

	ctx := waitForShutdown()

	// init database
	dbPool, err := postgres.New(cfg.PostgresConfig)
	if err != nil {
		log.Fatalf("Failed to initialize postgres: %v", err)
	}
	defer dbPool.Close()

	// init redis
	redisPool := redisInfra.NewRedisClient(cfg.RedisConfig)
	if err := redisPool.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize redis: %v", err)
	}
	defer redisPool.Close()

	// init minio
	minioStore, err := minio.New(cfg.MinioConfig)
	if err != nil {
		log.Fatalf("Failed to init minio: %v", err)
	}

	if err := minioStore.EnsureBucket(ctx); err != nil {
		log.Fatalf("Failed to ensure minio bucket: %v", err)
	}

	// init infras
	infraSession := sessionInfra.NewSessionRepo(dbPool.DB)

	// init middlewares
	authMiddleware := middleware.NewAuthMiddleware(infraSession, true)

	// init repos
	adminRepo := adminRepo.NewAdminRepo(dbPool.DB)
	authRepo := authRepo.NewAuthRepo(dbPool.DB)

	// init services
	hasher := security.NewBcryptHasher(10)
	redis := redisStore.NewOTPRedisStore(redisPool.Client)
	tokenSrv := security.NewToken(cfg.TokenConfig.AccessSecret, cfg.TokenConfig.RefreshSecret, 
		cfg.TokenConfig.AccessTTL, cfg.TokenConfig.RefreshTTL)


	// init usecasesp
	adminUsecase := adminUsecase.NewAdminUsecase(adminRepo, hasher)
	authUsecase := authUsecase.NewAuthUsecase(authRepo, infraSession, redis, tokenSrv, hasher)

	
	// init handler
	authHandler := auth.NewAuthHandler(authUsecase)
	print(adminUsecase)

	// init server
	srv := server.NewServer(cfg.Server, authMiddleware, authHandler)

	// start server async
	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// wait for signal
	<-ctx.Done()
	log.Println("Shutdown signal received")

	// shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	log.Println("Graceful shutdown completed")
}

func createSuperuser() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbPool, err := postgres.New(cfg.PostgresConfig)
	if err != nil {
		log.Fatalf("Failed to initialize postgres: %v", err)
	}
	defer dbPool.Close()

	repo := adminRepo.NewAdminRepo(dbPool.DB)
	hasher := security.NewBcryptHasher(8)

	usecase := adminUsecase.NewAdminUsecase(repo, hasher)

	err = usecase.CreateSuperuser()
	if err != nil {
		log.Fatalf("Failed to create superuser: %v", err)
	}
}
