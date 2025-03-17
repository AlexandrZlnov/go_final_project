package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AlexandrZlnov/go_final_project/models"
	"github.com/AlexandrZlnov/go_final_project/service"
)

// хэндлер обработчик GET-запроса, по адресу /api/task
// возвращает все параметры задачи по её id
// GET-запрос приходит в формате /api/task?id=<идентификатор>
// клиенту возвращается JSON-объект со всеми полями задачи
func GetEditTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var task models.Task

	//проверяем токен на валидность
	token, err := service.VerifyToken(r)
	if err != nil || !token.Valid {
		log.Println("Получили не валидный токен: в GetEditTask")
		service.Error(w, "Authentification required", http.StatusUnauthorized)
		return
	}

	taskID := r.FormValue("id")
	if taskID == "" {
		service.Error(w, "Не получен ID записи", http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT ID, Date, Title, Comment, Repeat FROM scheduler WHERE id = :id",
		sql.Named("id", taskID))
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		service.Error(w, "Ошибка сканирования БД", http.StatusInternalServerError)
		return
	}
	service.Success(w, task, http.StatusOK)
}
