package sessionInfra

import "database/sql"

type sessionRepo struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *sessionRepo {
	return &sessionRepo{
		db: db,
	}
}
