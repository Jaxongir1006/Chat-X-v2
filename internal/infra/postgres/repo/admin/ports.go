package admin

type AdminStore interface {
	CreateSuperuser(email, password string) error
}
