package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// func HashPassword(password string) (string, error) {
// 	bytes ,err := bcrypt.GenerateFromPassword([]byte(password), 14)

// 	return string(bytes), err
// }
func HashPassword(password string) (string, error) {
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedBytes), nil
}
