package model

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Balance  int    `json:"balance"`
}
