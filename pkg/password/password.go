package password

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	// Use cost 12 for better security (default is 10)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword compares a hashed password with a plain text password
func VerifyPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// IsValidPassword checks if a password meets minimum requirements
func IsValidPassword(password string) bool {
	// Minimum 6 characters, maximum 100 characters
	return len(password) >= 6 && len(password) <= 100
}
