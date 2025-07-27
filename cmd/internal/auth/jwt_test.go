package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateTestToken(secret, username string, expired bool) string {
	expirationTime := time.Now().Add(5 * time.Minute)
	if expired {
		expirationTime = time.Now().Add(-5 * time.Minute)
	}

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokStr, _ := token.SignedString([]byte(secret))
	return tokStr
}

func TestValidateToken_Valid(t *testing.T) {
	// Set the JWT_SECRET env variable used in configparser
	os.Setenv("JWT_SECRET", "my-secret")

	token := generateTestToken("my-secret", "testuser", false)
	claims, err := validateToken(token)

	assert.NoError(t, err)
	assert.Equal(t, "testuser", claims.Username)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "my-secret")

	token := "invalid.token.string"
	claims, err := validateToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "my-secret")

	token := generateTestToken("my-secret", "testuser", true)
	claims, err := validateToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "my-secret")

	token := generateTestToken("my-secret", "testuser", false)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	handler := JWTMiddleware(next)
	handler.ServeHTTP(rr, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestJWTMiddleware_MissingToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler := JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not call next")
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Missing or malformed token")
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.string")
	rr := httptest.NewRecorder()

	handler := JWTMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not call next")
	}))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid or expired token")
}
