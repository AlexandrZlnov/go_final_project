package models

import (
	"fmt"
	"time"

	"github.com/AlexandrZlnov/go_final_project/service"
)

type Task struct {
	ID      string `json:"id"`      // идентификатор задачи
	Date    string `json:"date"`    // дата задачи в формате 20060102
	Title   string `json:"title"`   // заголовок задачи
	Comment string `json:"comment"` // комментарий к задаче
	Repeat  string `json:"repeat"`  // правило повторения
}

type Pass struct {
	Password string `json:"password"`
}

func (t Task) ValidateTaskData() (Task, error) {
	now := time.Now()

	if len(t.Title) == 0 {
		return t, fmt.Errorf("заголовок задачи не может быть пустым")
	}

	if t.Date == "" {
		t.Date = now.Format(service.DateFormat)
	}

	taskDate, err := time.Parse(service.DateFormat, t.Date)
	if err != nil {
		return t, fmt.Errorf("неподдерживаемый формат даты")
	}

	if taskDate.Format(service.DateFormat) < now.Format(service.DateFormat) {
		if t.Repeat == "" {
			t.Date = now.Format(service.DateFormat)
		}

		if len(t.Repeat) == 1 && t.Repeat == "w" {
			return t, fmt.Errorf("неподдерживаемый формат правила повторения задачи")
		}

		if t.Repeat != "" {
			nextDate, err := service.NextDate(now, t.Date, t.Repeat)
			if err != nil {
				return t, fmt.Errorf("неподдерживаемый формат правила повторения задачи")
			}
			t.Date = nextDate
		}
	}

	return t, nil
}
