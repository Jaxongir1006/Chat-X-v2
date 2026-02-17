package userInfra

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type userRepo struct {
	db     *sql.DB
	logger zerolog.Logger
}

func NewUserRepo(db *sql.DB, logger zerolog.Logger) *userRepo {
	return &userRepo{
		db:     db,
		logger: logger,
	}
}
