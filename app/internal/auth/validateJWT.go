package auth

import (
    "errors"
    "github.com/google/uuid"
    "github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
    claims := &jwt.RegisteredClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(tokenSecret), nil
    })
    if err != nil {
        return uuid.Nil, err
    }
    if !token.Valid {
        return uuid.Nil, errors.New("invalid token")
    }
    userID, err := uuid.Parse(claims.Subject)
    if err != nil {
        return uuid.Nil, err
    }
    return userID, nil
}