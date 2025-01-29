package auth

import "golang.org/x/crypto/bcrypt"

// Hash password
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// Compares a bcrypt hashed password with the given password
func ComparePassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
