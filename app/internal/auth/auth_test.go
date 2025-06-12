package auth

import (
    "testing"
    "net/http"
    "errors"
    "time"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
)

func TestMakeAndValidateJWT(t *testing.T) {
	secret := "supersecret"
    userID := uuid.New()
    expiresIn := time.Minute

    token, err := MakeJWT(userID, secret, expiresIn)
    if err != nil {
        t.Fatalf("failed to create JWT: %v", err)
    }

    gotID, err := ValidateJWT(token, secret)
    if err != nil {
        t.Fatalf("failed to validate JWT: %v", err)
    }
    if gotID != userID {
        t.Errorf("expected userID %v, got %v", userID, gotID)
    }
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
    secret := "supersecret"
    userID := uuid.New()
    expiresIn := -time.Minute // already expired

    token, err := MakeJWT(userID, secret, expiresIn)
    if err != nil {
        t.Fatalf("failed to create JWT: %v", err)
    }

    _, err = ValidateJWT(token, secret)
    if err == nil {
        t.Error("expected error for expired token, got nil")
    }
}

func TestValidateJWT_WrongSecret(t *testing.T) {
    secret := "supersecret"
    wrongSecret := "nottherightsecret"
    userID := uuid.New()
    expiresIn := time.Minute

    token, err := MakeJWT(userID, secret, expiresIn)
    if err != nil {
        t.Fatalf("failed to create JWT: %v", err)
    }

    _, err = ValidateJWT(token, wrongSecret)
    if err == nil {
        t.Error("expected error for token signed with wrong secret, got nil")
    }
}

func TestGetBearerToken(t *testing.T) {
    tests := []struct {
        name    string
        header  string
        want    string
        wantErr bool
    }{
		{"valid token", "Bearer abc123", "abc123", false},
        {"missing header", "", "", true},
        {"wrong prefix", "Token abc123", "", true},
        {"empty token", "Bearer ", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            headers := http.Header{}
            if tt.header != "" {
                headers.Set("Authorization", tt.header)
            }
            got, err := GetBearerToken(headers)
            if (err != nil) != tt.wantErr {
                t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
            }
            if got != tt.want {
                t.Errorf("expected token: %q, got: %q", tt.want, got)
            }
        })
    }
}