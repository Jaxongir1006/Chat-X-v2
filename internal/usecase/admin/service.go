package admin

import (
	adminRepo "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/admin"
	"github.com/Jaxongir1006/Chat-X-v2/internal/infra/security"
)

type AdminUsecase struct {
	adminRepo adminRepo.AdminStore
	hasher    security.Hasher
}

func NewAdminUsecase(adminRepo adminRepo.AdminStore, hasher security.Hasher) *AdminUsecase {
	return &AdminUsecase{
		adminRepo: adminRepo,
		hasher:    hasher,
	}
}
