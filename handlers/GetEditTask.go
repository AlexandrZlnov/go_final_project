package handlers

import (
	"database/sql"
	"net/http"

	"github.com/AlexandrZlnov/go_final_project/models"
	"github.com/AlexandrZlnov/go_final_project/service"
)

func GetEditTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var task models.Task
	//task := models.Task{}

	taskID := r.FormValue("id")
	if taskID == "" {
		service.Error(w, "Не получен ID записи", http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id",
		sql.Named("id", taskID))
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		service.Error(w, "Ошибка сканирования БД", http.StatusInternalServerError)
		return
	}
	service.Success(w, task, http.StatusOK)
}
