package auth

import (
    "errors"
    "time"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
)

// HashPassword hashes the given password using bcrypt.
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPasswordHash compares a bcrypt hashed password with its possible plaintext equivalent.
// Returns nil on success, or an error on failure.
func CheckPasswordHash(hash, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

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