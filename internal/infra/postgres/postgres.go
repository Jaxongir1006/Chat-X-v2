package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/config"
	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

func New(cfg config.PostgresConfig) (*Postgres, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DB, cfg.SSL,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &Postgres{DB: db}, nil
}

func (p *Postgres) Close() error {
	if p.DB == nil {
		return nil
	}
	return p.DB.Close()
}

func (p *Postgres) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return p.DB.PingContext(ctx)
}
