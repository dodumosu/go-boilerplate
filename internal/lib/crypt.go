package lib

import (
	"github.com/go-crypt/crypt"
	"github.com/go-crypt/crypt/algorithm/argon2"
)

func VerifyPassword(plaintext string, hashedPassword string) (bool, error) {
	return crypt.CheckPassword(plaintext, hashedPassword)
}

func HashPassword(plaintext string) (string, error) {
	hasher, err := argon2.New(argon2.WithProfileRFC9106LowMemory())
	if err != nil {
		return "", err
	}
	digest, err := hasher.Hash(plaintext)
	if err != nil {
		return "", err
	}
	return digest.Encode(), nil
}
