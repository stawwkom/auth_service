package app

import (
	"context"
	"fmt"
	auth2 "github.com/stawwkom/auth_service/internal/api/auth"
	"github.com/stawwkom/auth_service/internal/config"
	a2 "github.com/stawwkom/auth_service/internal/delivery/grpc/access_handler"
	a1 "github.com/stawwkom/auth_service/internal/delivery/grpc/auth_handler"
	"github.com/stawwkom/auth_service/internal/service"
	"github.com/stawwkom/auth_service/internal/service/jwt_service"
	acServ "github.com/stawwkom/auth_service/internal/service/jwt_service/access"
	aServ "github.com/stawwkom/auth_service/internal/service/jwt_service/auth"
	descAccess "github.com/stawwkom/auth_service/pkg/access_v1"
	descAuth "github.com/stawwkom/auth_service/pkg/auth_login"
)

type serviceProvider struct {
	authService   service.AuthService
	dbClose       func() error
	aService      jwt_service.AuthService
	accessService jwt_service.AccessService
	authImpl      *auth2.Server
	accessImpl    descAccess.AccessV1Server
	aImpl         descAuth.AuthV1Server
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

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth2.Server {
	if s.authImpl == nil {
		s.authImpl = auth2.NewServer(s.AuthService())
	}
	return s.authImpl
}
func (s *serviceProvider) AuthS(ctx context.Context) jwt_service.AuthService {
	if s.aService == nil {
		s.aService = aServ.NewAuthService()
	}
	return s.aService
}
func (s *serviceProvider) AuthI(ctx context.Context) descAuth.AuthV1Server {
	if s.aImpl == nil {
		s.aImpl = a1.NewAuthHandler(s.AuthS(ctx))
	}
	return s.aImpl
}

func (s *serviceProvider) AccessS(ctx context.Context) jwt_service.AccessService {
	if s.accessService == nil {
		s.accessService = acServ.NewAccessService()
	}
	return s.accessService
}

func (s *serviceProvider) AccessI(ctx context.Context) descAccess.AccessV1Server {
	if s.accessImpl == nil {
		s.accessImpl = a2.NewAccessHandler(s.AccessS(ctx))
	}
	return s.accessImpl
}
