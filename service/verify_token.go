package service

import (
	"log"
	"net/http"

	"github.com/AlexandrZlnov/go_final_project/config"
	"github.com/golang-jwt/jwt/v5"
)

// верификация JWT-токена
func VerifyToken(r *http.Request) (*jwt.Token, error) {

	jwtSecret := config.GetJWTKey()

	// извлекаем куку с токеном
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Println("Ошибка извлечения токена из куки")
		return nil, err
	}

	// парсим токен
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("Проблема определения метода подписания токена:")
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
