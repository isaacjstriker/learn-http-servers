package auth

import (
    "time"
    "github.com/google/uuid"
    "github.com/golang-jwt/jwt/v5"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
    now := time.Now().UTC()
    claims := jwt.RegisteredClaims {
        Issuer: "chirpy",
        IssuedAt: jwt.NewNumericDate(now),
        ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
        Subject: userID.String(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString([]byte(tokenSecret))
    if err != nil {
        return "", err
    }

    return signedToken, nil
}