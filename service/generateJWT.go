// генерируем и подписываем ключем токен
package service

import (
	"time"

	"github.com/AlexandrZlnov/go_final_project/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken() (string, error) {

	jwtSecret := config.GetJWTKey()

	// формируем Claims
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	// создаем новый токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// подписываем токен с ключем
	TokenSign, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return TokenSign, nil
}
