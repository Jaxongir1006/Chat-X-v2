package adminRepo

import "database/sql"

type adminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) *adminRepo {
	return &adminRepo{
		db: db,
	}
}
