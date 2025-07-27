package auth

import (
	configparser "companies/cmd/internal/configParser"
	"companies/cmd/internal/consts"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// var jwtKey = []byte("your-very-secret-key") // replace with env/config

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func validateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	jwtKey := configparser.GetCfgValue("JWT_SECRET", "default-secret")
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(consts.ApplicationPrefix, "JWT middleware")

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or malformed token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := validateToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
