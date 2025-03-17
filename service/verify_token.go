package service

import (
	"log"
	"net/http"

	"github.com/AlexandrZlnov/go_final_project/config"
	"github.com/golang-jwt/jwt/v5"
)

//var jwtSecret = []byte("mysecretkey")

func VerifyToken(r *http.Request) (*jwt.Token, error) {

	jwtSecret := config.GetJWTKey()
	log.Println("Секретный ключ из VerifyToken ------>", jwtSecret, "   ", string(jwtSecret))

	// извлекаем куку с токеном
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Println("Ошибка извлечения токена из куки")
		return nil, err
	}
	log.Println("Токен из куки:", cookie)

	// парсим токен
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		log.Println("Распарсили токен:", token)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("Проблема определения метода подписания токена:")
			return nil, jwt.ErrSignatureInvalid
		}
		log.Println("Возвращаем секретный ключ после парсинга", jwtSecret)
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	log.Println("----------------------")
	return token, nil
}
