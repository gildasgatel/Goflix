package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPasswd(pswd []byte) ([]byte, error) {
	hashMotDePasse, err := bcrypt.GenerateFromPassword(pswd, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashMotDePasse, nil
}

func CompareHashAndPassword(pswd, hashPswd []byte) error {
	return bcrypt.CompareHashAndPassword(hashPswd, pswd)

}
