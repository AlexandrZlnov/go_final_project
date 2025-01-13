package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/AlexandrZlnov/go_final_project/models"
	"github.com/AlexandrZlnov/go_final_project/service"
)

//var dateFormat = "20060102"

func PostAddTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var newTask models.Task
	var bufer bytes.Buffer

	_, err := bufer.ReadFrom(r.Body)
	if err != nil {
		service.Error(w, "ReadFrom error", http.StatusBadRequest)
		return
	}

	// if err = json.Unmarshal(bufer.Bytes(), &newTask); err != nil {
	// 	service.Error(w, "Unmarshal error", http.StatusBadRequest)
	// 	return
	err = json.Unmarshal(bufer.Bytes(), &newTask)
	if err != nil {
		service.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	validTask, err := models.Task.ValidateTaskData(newTask)
	if err != nil {
		// service.Error(w, "Error validating task", http.StatusBadRequest)
		// return
		service.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTask = validTask

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", newTask.Date),
		sql.Named("title", newTask.Title),
		sql.Named("comment", newTask.Comment),
		sql.Named("repeat", newTask.Repeat))

	if err != nil {
		service.Error(w, "Error inserting task in to DB4", http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		service.Error(w, "Error getting last insert id", http.StatusInternalServerError)
		return
	}

	//service.Success(w, id, http.StatusOK)
	jsonResponse := map[string]int64{"id": id}
	service.Success(w, jsonResponse, http.StatusCreated)
}
