package jwtutils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/stawwkom/auth_service/internal/model"
	"time"
)

func GenerateToken(info model.UserInformation, secretKey []byte, duration time.Duration) (string, error) {
	claim := model.UserClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Username: info.Username,
		Role:     info.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(secretKey)
}

func VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &model.UserClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected token signing method")
			}
			return secretKey, nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("invalid token: %v", err.Error())
	}

	claim, ok := token.Claims.(*model.UserClaim)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claim, nil
}
