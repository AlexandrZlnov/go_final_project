package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/AlexandrZlnov/go_final_project/models"
	"github.com/AlexandrZlnov/go_final_project/service"
)

func PostDoneTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	task := models.Task{}
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

	if task.Date == "" {
		task.Date = time.Now().Format(dateFormat)
	}

	if task.Repeat == "" {
		_, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
			sql.Named("id", task.ID))
		if err != nil {
			service.Error(w, "Ошибка удаления записи из БД", http.StatusInternalServerError)
			return
		}
	} else if task.Repeat != "" {
		nextDate, err := service.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			service.Error(w, "Ошибка расчета следующей даты задачи", http.StatusInternalServerError)
			return
		}
		task.Date = nextDate
		_, err = db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
			sql.Named("date", task.Date),
			sql.Named("id", task.ID))
		if err != nil {
			service.Error(w, "Ошибка обновления даты задачи в БД", http.StatusInternalServerError)
			return
		}
	}

	jsonResponse := map[string]models.Task{}
	service.Success(w, jsonResponse, http.StatusOK)
}
