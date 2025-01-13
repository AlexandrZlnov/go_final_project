package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/AlexandrZlnov/go_final_project/models"
	"github.com/AlexandrZlnov/go_final_project/service"
)

func PutEditTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var changedTask models.Task
	var bufer bytes.Buffer

	_, err := bufer.ReadFrom(r.Body)
	if err != nil {
		service.Error(w, "ReadFrom error", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bufer.Bytes(), &changedTask)
	if err != nil {
		service.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	validTask, err := models.Task.ValidateTaskData(changedTask)
	if err != nil {
		service.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	changedTask = validTask

	res, err := db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("id", changedTask.ID),
		sql.Named("date", changedTask.Date),
		sql.Named("title", changedTask.Title),
		sql.Named("comment", changedTask.Comment),
		sql.Named("repeat", changedTask.Repeat))

	if err != nil {
		service.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	row, err := res.RowsAffected()
	if err != nil {
		service.Error(w, "Error getting rows affected", http.StatusInternalServerError)
		return
	}

	if row == 0 {
		service.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	jsonResponse := map[string][]models.Task{}
	service.Success(w, jsonResponse, http.StatusOK)

}
