package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AlexandrZlnov/go_final_project/models"
	"github.com/AlexandrZlnov/go_final_project/service"
)

// константа определяет максимальное значение выводимых записей
const taskLimitPerPage = 15

// хэндлер обработчик POST запроса по адресу /api/tasks
// возвращает список ближайших задач в формате JSON в виде списка в поле tasks
// задачи отсортированы по дате в сторону увеличенияя
// каждая задача содержит все поля таблицы scheduler в виде строк
// дата представлена в формате 20060102

// так же возвращает список задач по поисковому запросу
// полученному в строке запроса в формате /api/tasks?search=поисковый запрос

func GetTasks(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	token, err := service.VerifyToken(r)
	if err != nil || !token.Valid {
		log.Println("Получили не валидный токен")
		service.Error(w, "Authentification required", http.StatusUnauthorized)
		return
	}

	var tasks []models.Task

	searchQuery := r.FormValue("search")

	if len(searchQuery) > 0 {
		dateTime, err := time.Parse("02.01.2006", searchQuery)
		if err != nil {
			rows, err := db.Query("SELECT ID, Date, Title, Comment, Repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit ",
				sql.Named("search", fmt.Sprint("%"+searchQuery+"%")), sql.Named("limit", taskLimitPerPage))
			if err != nil {
				service.Error(w, "Ошибка запроса к БД по полю title", http.StatusInternalServerError)
				return
			}
			tasks = processingDBQueryResults(w, rows)
		} else {
			rows, err := db.Query("SELECT ID, Date, Title, Comment, Repeat FROM scheduler WHERE date = :date ORDER BY date LIMIT :limit",
				sql.Named("date", dateTime.Format(service.DateFormat)), sql.Named("limit", taskLimitPerPage))
			if err != nil {
				service.Error(w, "Ошибка запроса БД по полю date", http.StatusInternalServerError)
				return
			}
			tasks = processingDBQueryResults(w, rows)
		}
	} else {
		rows, err := db.Query("SELECT ID, Date, Title, Comment, Repeat FROM scheduler ORDER BY date LIMIT :limit",
			sql.Named("limit", taskLimitPerPage))
		if err != nil {
			service.Error(w, "Ошибка запроса БД", http.StatusInternalServerError)
			return
		}
		tasks = processingDBQueryResults(w, rows)
	}

	jsonResponse := map[string][]models.Task{"tasks": tasks}
	service.Success(w, jsonResponse, http.StatusOK)
}

func processingDBQueryResults(w http.ResponseWriter, rows *sql.Rows) []models.Task {
	var tasks []models.Task

	defer rows.Close()

	for rows.Next() {
		task := models.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			service.Error(w, "Ошибка сканирования БД", http.StatusInternalServerError)
			return nil
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		service.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return nil
	}

	if len(tasks) == 0 {
		tasks = make([]models.Task, 0)
	}

	if err := rows.Err(); err != nil {
		service.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return nil
	}

	return tasks
}
