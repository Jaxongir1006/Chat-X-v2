package adminUsecase

type AdminUsecase struct {
	adminRepo AdminRepo
	hasher    Hasher
}

func NewAdminUsecase(adminRepo AdminRepo, hasher Hasher) *AdminUsecase {
	return &AdminUsecase{
		adminRepo: adminRepo,
		hasher:    hasher,
	}
}
