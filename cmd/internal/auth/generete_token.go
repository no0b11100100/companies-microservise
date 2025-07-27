package auth

import (
	configparser "companies/cmd/internal/configParser"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtKey := configparser.GetCfgValue("JWT_SECRET", "very-secret-key")
	return token.SignedString([]byte(jwtKey))
}

func HandleFunc(w http.ResponseWriter, r *http.Request) {
	token, err := GenerateToken("admin")
	if err != nil {
		log.Println("TOKEN ERROR", err)
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}
