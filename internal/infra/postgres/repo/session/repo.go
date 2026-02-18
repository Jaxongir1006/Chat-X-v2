package sessionInfra

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog"
)

type sessionRepo struct {
	db     *sql.DB
	tx     *sql.Tx
	logger zerolog.Logger
}

func NewSessionRepo(db *sql.DB, logger zerolog.Logger) *sessionRepo {
	return &sessionRepo{
		db:		db,
		logger: logger,
	}
}

func (r *sessionRepo) WithTx(tx *sql.Tx) *sessionRepo {
	return &sessionRepo{db: r.db, tx: tx, logger: r.logger}
}

func (r *sessionRepo) execer() interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
} {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

