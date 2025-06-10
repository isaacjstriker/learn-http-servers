package auth

import (
    "golang.org/x/crypto/bcrypt"
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