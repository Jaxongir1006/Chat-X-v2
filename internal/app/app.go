package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/config"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/minio"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres"
	adminRepo "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/admin"
	redisInfra "github.com/Jaxongir1006/Chat-X-v2/internal/infra/redis"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/security"
	"github.com/Jaxongir1006/Chat-X-v2/internal/server"
	"github.com/Jaxongir1006/Chat-X-v2/internal/usecase/adminUsecase"
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

	// init
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

	srv := server.NewServer(cfg.Server)

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
