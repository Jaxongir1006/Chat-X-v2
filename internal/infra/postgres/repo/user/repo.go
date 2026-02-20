package user

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog"
)

type userRepo struct {
	db     *sql.DB
	tx     *sql.Tx
	logger zerolog.Logger
}

func NewUserRepo(db *sql.DB, logger zerolog.Logger) *userRepo {
	return &userRepo{
		db:     db,
		logger: logger,
	}
}

func (r *userRepo) WithTx(tx *sql.Tx) *userRepo {
	return &userRepo{db: r.db, tx: tx, logger: r.logger}
}

func (r *userRepo) execer() interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
} {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}
