package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Role  string `json:"role"`
	Email string `json:"email"`
	jwt.StandardClaims
}
