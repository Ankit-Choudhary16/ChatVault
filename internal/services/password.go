package services

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword creates a bcrypt hash from the given password.
// Returns the hashed password as a string, or an error if hashing fails.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ComparePassword checks if the given password matches the stored hash.
// Returns true if the password is valid, false otherwise.
func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
