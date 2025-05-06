package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)



func encryptPassword(pass string) string{
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	if err != nil{
		log.Fatal("Error while hashing password")
	}
	return string(hash)
}

func comparePassword(pass, hash string)bool{
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(pass))

	if err != nil{
		return false
	}
	return true
}