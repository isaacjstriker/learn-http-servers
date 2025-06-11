package auth

import (
    "testing"
    "time"

    "github.com/google/uuid"
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