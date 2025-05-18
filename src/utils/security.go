package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrInternal.Wrap(err, "hash password")
	}
	return string(hashBytes), nil
}

func GenerateCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid code length: %d", length)
	}
	var result []byte
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		result = append(result, byte('0'+num.Int64()))
	}
	return string(result), nil
}
