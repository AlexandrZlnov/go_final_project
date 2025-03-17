package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AlexandrZlnov/go_final_project/models"
	"github.com/AlexandrZlnov/go_final_project/service"
)

// хэндлер обработчик DELETE запроса по адресу /api/pask
// удаляет задачу по id
// запроса методом DELETE приходит от клиента приходит в формате: /api/task/done?id=<идентификатор>.
func DeleteTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	token, err := service.VerifyToken(r)
	if err != nil || !token.Valid {
		log.Println("Получили не валидный токен")
		service.Error(w, "Authentification required", http.StatusUnauthorized)
		return
	}

	taskID := r.FormValue("id")

	if taskID == "" {
		service.Error(w, "Не получен ID записи", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", taskID))
	if err != nil {
		service.Error(w, "Ошибка удаления задачи из БД", http.StatusInternalServerError)
		return
	}

	row, err := res.RowsAffected()
	if err != nil {
		service.Error(w, "Ошибка обновления db", http.StatusInternalServerError)
		return
	}

	if row == 0 {
		service.Error(w, "Задача не найдена", http.StatusInternalServerError)
		return
	}

	jsonResponse := map[string]models.Task{}
	service.Success(w, jsonResponse, http.StatusOK)

}
