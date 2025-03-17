package config

import (
	"log"
	"os"
)

func GetJWTKey() []byte {
	secretKeyStr, ok := os.LookupEnv("TODO_SCRTKEY")
	if !ok || secretKeyStr == "" {
		log.Fatal("ERROR: TODO_SCRTKEY environment variable is not set")
	}
	JWTKey := []byte(secretKeyStr)
	if len(JWTKey) == 0 {
		log.Fatal("ERROR: secretKey variable is empty")
	}
	return JWTKey
}

func GetUserPass() string {
	password, ok := os.LookupEnv("TODO_PASSWORD")
	if !ok || password == "" {
		log.Fatal("ERROR: TODO_PASSWORD environment variable is not set")
	}
	// password := []byte(passwordStr)
	// if len(password) == 0 {
	// 	log.Fatal("ERROR: password variable is empty")
	// }
	return password
}
