package authRepo

import "database/sql"

type authRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *authRepo {
	return &authRepo{
		db: db,
	}
}
