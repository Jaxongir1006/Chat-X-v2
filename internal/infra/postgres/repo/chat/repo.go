package chat

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog"
)

type chatRepo struct {
	db     *sql.DB
	tx     *sql.Tx
	logger zerolog.Logger
}

func NewChatRepo(db *sql.DB, logger zerolog.Logger) *chatRepo {
	return &chatRepo{
		db:     db,
		logger: logger,
	}
}

func (r *chatRepo) WithTx(tx *sql.Tx) *chatRepo {
	return &chatRepo{db: r.db, tx: tx, logger: r.logger}
}

func (r *chatRepo) execer() interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
} {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}
