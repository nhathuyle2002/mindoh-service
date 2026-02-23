package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *AuthService) GenerateJWT(username string, role Role) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))
}

func (s *AuthService) ParseAndValidateJWT(tokenString string) (string, Role, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return "", "", err
	}
	return claims.Username, claims.Role, nil
}
