package auth

import "github.com/Jaxongir1006/Chat-X-v2/internal/config"

type authHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *authHandler {
	return &authHandler{
		cfg: cfg,
	}
}
