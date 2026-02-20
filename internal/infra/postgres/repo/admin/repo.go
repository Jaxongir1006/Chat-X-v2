package admin

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type adminRepo struct {
	db     *sql.DB
	logger zerolog.Logger
}

func NewAdminRepo(db *sql.DB, logger zerolog.Logger) *adminRepo {
	return &adminRepo{
		db:     db,
		logger: logger,
	}
}
