package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type Claims struct {
	UserID uint `json:"user_id"`
	Role   Role `json:"role"`
	jwt.RegisteredClaims
}

type AuthContext struct {
	UserID uint `json:"user_id"`
	Role   Role `json:"role"`
}
