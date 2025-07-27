package integration

import (
	"companies/cmd/internal/auth"
	configparser "companies/cmd/internal/configParser"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &auth.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtKey := configparser.GetCfgValue("JWT_SECRET", "default-secret")
	return token.SignedString(jwtKey)
}
