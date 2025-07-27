//TODO: remove this file

package auth

import (
	configparser "companies/cmd/internal/configParser"
	"encoding/json"
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

	jwtKey := configparser.GetCfgValue("JWT_SECRET", "default-secret")
	return token.SignedString(jwtKey)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleFunc(w http.ResponseWriter, r *http.Request) {
	var creds LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Dummy check — replace with DB/user service check
	if creds.Username != "admin" || creds.Password != "password" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(creds.Username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
