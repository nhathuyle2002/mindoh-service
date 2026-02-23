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
	Username string `json:"username"`
	Role     Role   `json:"role"`
	jwt.RegisteredClaims
}

type AuthContext struct {
	UserID   uint
	Username string
	Role     Role
}
