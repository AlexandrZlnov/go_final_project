package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AlexandrZlnov/go_final_project/service"
)

// хэндлер обработчик GET-запроса, по адресу api/nextdate
// GET-запрос приходит в формате: "/api/nextdate?now=<20060102>&date=<20060102>&repeat=<правило>"
// обработчик вызвает функцию NextDate из пакета service и возвращает дату следующего выполнения задачи в формате 20060102
func GetNextDate(w http.ResponseWriter, r *http.Request) {
	date := r.FormValue("date")

	if date == "" {
		http.Error(w, "Не указана дата", http.StatusBadRequest)
		return
	}

	repeat := r.FormValue("repeat")
	now, err := time.Parse(service.DateFormat, r.FormValue("now"))
	if err != nil {
		http.Error(w, "Неверный формат даты now", http.StatusBadRequest)
		return
	}

	nextDate, err := service.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, "Неудачная попытка вычисления даты переноса", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, nextDate)
}
