package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	//"os"

	"github.com/AlexandrZlnov/go_final_project/config"
	"github.com/AlexandrZlnov/go_final_project/models"
	"github.com/AlexandrZlnov/go_final_project/service"
)

// хэндлер обработчик POST запроса регистрации пользователя /api/signin
// форма содержить только поле для ввода пароля, поэтому именит пользователя нет
// обработчик получает JSON с полем password, которое содержит введённый пароль
// функция сверяет указанный пароль с хранимым в переменной окружения TODO_PASSWORD
// если они совпадают, формируется JWT-токен и возвращается клиенту поле token JSON-объекта
// если пароль неверный или произошла ошибка, возвращается JSON c текстом ошибки в поле error
func PostSignin(w http.ResponseWriter, r *http.Request) {
	//pass := os.Getenv("TODO_PASSWORD")
	password := config.GetUserPass()

	log.Println("Пароль в переменной ENV:", password)

	var passIn models.Pass
	var bufer bytes.Buffer

	_, err := bufer.ReadFrom(r.Body)
	if err != nil {
		service.Error(w, "ReadFrom error", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bufer.Bytes(), &passIn)
	if err != nil {
		service.Error(w, "Error decoding JSON", http.StatusBadRequest)
		log.Println("Error decoding JSON - read user password")
		return
	}

	log.Println("Входящий пароль:", passIn.Password)

	if passIn.Password != password {
		log.Println("Incorrect password")
		//service.Error(w, "Error decoding JSON", http.StatusBadRequest)

		response := map[string]string{
			"error": "Неверный пароль",
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Println("Ошибка кодирования ответа с токеном: Неверный пароль")
			return
		}

	} else {
		log.Println("Correct password")

		// получаем подписанный токен
		token, err := service.GenerateToken()
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			log.Println("Ошибка генарации токена")
		}

		// формируем JSON ответ с токеном
		response := map[string]string{
			"token": token,
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Println("Ошибка кодирования ответа с токеном")
			return
		}
	}

}
