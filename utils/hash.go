package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (hash string, err error) {
	var hashBytes []byte
	hashBytes, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashBytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
