package admin

func (a *adminRepo) CreateSuperuser(email, password string) error {
	query := `INSERT INTO users (email, password_hash, role, verified, created_at, updated_at) VALUES ($1, $2, 'superuser', TRUE, NOW(), NOW())`

	_, err := a.db.Exec(query, email, password)
	if err != nil {
		return err
	}

	return nil
}
