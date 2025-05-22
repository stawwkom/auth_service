package access_handler

import (
	"context"
	"errors"
	"github.com/stawwkom/auth_service/internal/service/jwt_service"
	descAccess "github.com/stawwkom/auth_service/pkg/access_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

const (
	authPrefix = "Bearer "
)

var accessibleRoles map[string]string

type accessServer struct {
	descAccess.UnimplementedAccessV1Server
	accessService jwt_service.AccessService
}

func NewAccessHandler(service jwt_service.AccessService) descAccess.AccessV1Server {
	return &accessServer{
		accessService: service,
	}
}

func (a *accessServer) Check(ctx context.Context, req *descAccess.CheckRequest) (*emptypb.Empty, error) {
	mb, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	authHeder, ok := mb["authorization"]
	if !ok || len(authHeder) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeder[0], authPrefix) {
		return nil, errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeder[0], authPrefix)

	err := a.accessService.CheckAccess(ctx, accessToken, req.GetEndpointAddress())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
