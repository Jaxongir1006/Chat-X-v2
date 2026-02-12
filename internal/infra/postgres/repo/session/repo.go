package sessionInfra

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type sessionRepo struct {
	db     *sql.DB
	logger zerolog.Logger
}

func NewSessionRepo(db *sql.DB, logger zerolog.Logger) *sessionRepo {
	return &sessionRepo{
		db:     db,
		logger: logger,
	}
}
