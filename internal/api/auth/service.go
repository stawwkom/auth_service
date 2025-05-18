package auth

import (
	"github.com/stawwkom/auth_service/internal/service"
	desc "github.com/stawwkom/auth_service/pkg/auth_v1"
)

// Server Реализует gRPC server
type Server struct {
	desc.UnimplementedUserAPIServer
	authService service.AuthService
}

func NewServer(authService service.AuthService) *Server {
	return &Server{
		authService: authService,
	}
}
