package helpers

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HassPass(pass string) (string, error) {
	if len(pass) == 0 {
		return "", errors.New("Password should not be empty")
	}

	bytePassword := []byte(pass)
	hashPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.MinCost)
	if err != nil {
		fmt.Printf("[UserController.SetPassword] error when generate password with error: %v\n", err)
		return "", nil
	}

	return string(hashPassword), nil
}