package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPass(password string) (string, error) {
	salt := make([]byte, 16)

	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("error generate salt %w", err)
	}

	plainPassword := []byte(password + string(salt))

	passBcrypt, _ := bcrypt.GenerateFromPassword(plainPassword, bcrypt.DefaultCost)

	saltAndPassword := append(salt, passBcrypt...)
	return base64.StdEncoding.EncodeToString(saltAndPassword), nil
}

func HashPassValidator(password, hash string) bool {
	saltAndPassword, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false
	}

	salt := saltAndPassword[:16]
	hashedPassword := saltAndPassword[16:]

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password+string(salt)))
	if err != nil {
		return false
	}

	return true
}
