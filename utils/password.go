//		Copyright (c) Itoplist - All Rights Reserved
//
//	Unauthorized copying of this file, via any medium is strictly prohibited
//	Proprietary and confidential
//	Written by Ilyes Cherfaoui <ogfris@protonmail.com>, 2019

package utils

import "golang.org/x/crypto/bcrypt"

// Encrypt encrypts passwords and return the hash with an error if there's one.
func Encrypt(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(passwordBytes), err
}

// Compare checks if a hashed password matches the password provided.
func Compare(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
