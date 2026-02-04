package adminRepo


type AdminStore interface {
	CreateSuperuser(email, password string) error
}