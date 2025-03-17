package service

import (
	"log"
	//"os"
	"time"

	"github.com/AlexandrZlnov/go_final_project/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken() (string, error) {

	jwtSecret := config.GetJWTKey()
	log.Println("Секретный ключ ------>", jwtSecret, "   ", string(jwtSecret))

	// формируем Claims
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	// создаем новый токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println("Создан токен:", token)

	// подписываем токен с ключем
	TokenSign, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("Ошибка подписания токена", TokenSign)
		return "", err
	}
	log.Println("Токер подписан:", TokenSign)
	return TokenSign, nil
}
