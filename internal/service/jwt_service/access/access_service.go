package access

import (
	"context"
	"errors"
	jwtu "github.com/stawwkom/auth_service/internal/jwtutils"
	"github.com/stawwkom/auth_service/internal/model"
	"github.com/stawwkom/auth_service/internal/service/jwt_service"
)

const (
	accessTokenSecretKey = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="
)

type accessServiceImpl struct {
	rolePermissions map[string]string
}

func NewAccessService() jwt_service.AccessService {
	return &accessServiceImpl{
		rolePermissions: map[string]string{
			model.ExamplePath: "admin",
		},
	}
}

func (a *accessServiceImpl) CheckAccess(ctx context.Context, token string, endpoint string) error {
	claims, err := jwtu.VerifyToken(token, []byte(accessTokenSecretKey))
	if err != nil {
		return errors.New("invalid access token")
	}
	requiredRole, ok := a.rolePermissions[endpoint]
	if !ok {
		// Если endpoint не зарегистрирован — доступ открыт
		return nil
	}
	if claims.Role != requiredRole {
		return errors.New("access denied: insufficient role")
	}
	return nil
}
