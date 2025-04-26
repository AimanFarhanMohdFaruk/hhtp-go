package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	pass := []byte(password)
	hashPassword, err := bcrypt.GenerateFromPassword(pass, 10)

	if err != nil {
		return err.Error(), err
	}

	return string(hashPassword), nil
}

func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	
	if err != nil {
		return err
	}

	return nil
}