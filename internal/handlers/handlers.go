package handlers

import "github.com/golates/api-gateway/pkg/config"

type Handlers struct {
	cfg         *config.Config
	AuthHandler *AuthHandler
}

func NewHandlers(cfg *config.Config) *Handlers {
	return &Handlers{
		cfg:         cfg,
		AuthHandler: NewAuthHandler(cfg.AuthConfig.URL),
	}
}
