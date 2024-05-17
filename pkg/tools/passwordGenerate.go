package tools

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(p string)string{
	password, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Gagal Generate Password")
	}
	return string(password)
}


