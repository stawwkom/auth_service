package app

import (
	"fmt"
	"github.com/stawwkom/auth_service/internal/config"
	"github.com/stawwkom/auth_service/internal/service"
)

type serviceProvider struct {
	authService service.AuthService
	dbClose     func() error
}

func (s *serviceProvider) AuthService() service.AuthService {
	return s.authService
}

func (s *serviceProvider) GRPCAddr() string {
	return config.Cfg.Server.Host + ":" + fmt.Sprintf("%d", config.Cfg.Server.Port)
}

func (s *serviceProvider) HTTPAddr() string {
	Port := config.Cfg.Server.Port + 1
	return config.Cfg.Server.Host + ":" + fmt.Sprintf("%d", Port)
}
