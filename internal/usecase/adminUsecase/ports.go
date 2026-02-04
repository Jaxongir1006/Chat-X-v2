package adminUsecase

type AdminRepo interface {
	CreateSuperuser(email, password string) error
}

type Hasher interface {
	Hash(password string) (string, error)
	CheckPasswordHash(password, hash string) error
}
