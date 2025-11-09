package jwt

import (
	"sso/internal/domain/models"
	"testing"
	"time"
)

func TestNewToken_ValidJWTWithCorrectClaims(t *testing.T) {
	user := models.User{
		ID:       42,
		Email:    "user@test.com",
		PassHash: []byte("fake-hash"),
	}

	app := models.App{
		ID:     7,
		Name:   "TestApp",
		Secret: "super-secret",
	}

	ttl := time.Hour

	tokenString, err := NewToken(user, app, ttl)
	if err != nil {
		t.Fatalf("expected no error from NewToken, got: %v", err)
	}
	if tokenString == "" {
		t.Fatal("expected non-empty token string")
	}
}
