package model

import "github.com/dgrijalva/jwt-go"

const (
	ExamplePath = "/note_v1.NoteV1/Get"
)

type UserClaim struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserInformation struct {
	Username string
	Role     string
}
