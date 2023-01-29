package entity

import "github.com/golang-jwt/jwt/v4"

type UserClaims struct {
	User   *User       `json:"user"`
	Device *UserDevice `json:"device"`
	jwt.RegisteredClaims
}
