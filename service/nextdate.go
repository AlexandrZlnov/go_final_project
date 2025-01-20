package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	DateFormat      = "20060102"
	TransferDaysMax = 400
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" || repeat == " " {
		return "", fmt.Errorf("правило повторения не должно содержать пустую строку, %s", repeat)
	}

	nextDate, err := time.Parse(DateFormat, date)
	if err != nil {
		return fmt.Sprintf("неверный формат даты: %s", date), err
	}

	// Определеяем правило переноса задач
	transferRule := strings.Split(repeat, "")[0]

	if !(strings.EqualFold(transferRule, "d") || strings.EqualFold(transferRule, "y") || strings.EqualFold(transferRule, "w")) {
		return "", fmt.Errorf("неподдерживаемый формат повторения задачи")
	}

	switch transferRule {

	case "d":
		if len(repeat) < 2 {
			return "", fmt.Errorf("неподдерживаемый формат правила повторения даты")
		}
		transferDays := strings.Split(repeat, " ")[1:]
		transferDaysInt, err := strconv.Atoi(transferDays[0])
		if err != nil {
			return "", fmt.Errorf("не верный формат количества дней переноса заметки")
		}

		if transferDaysInt < 1 || transferDaysInt > 400 {
			return "", fmt.Errorf("количество дней для переноса должно быть от 1 до %d", TransferDaysMax)
		}

		nextDate = nextDate.AddDate(0, 0, transferDaysInt)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, transferDaysInt)
		}

	case "y":
		nextDate = nextDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}

	case "w":
		splitedRepeatRule := strings.Split(repeat, " ")[1]
		if len(splitedRepeatRule) == 0 {
			return "", fmt.Errorf("не указан параметр дней недели для переноса заметки: %s", repeat)
		}

		transferWeekdayByDay := strings.Split(splitedRepeatRule, ",")

		weekDaysInt := make([]int, 0, len(transferWeekdayByDay))

		for _, weekDay := range transferWeekdayByDay {
			weekDayInt, err := strconv.Atoi(weekDay)
			if err != nil {
				return "", fmt.Errorf("не верный формат дней недели для переноса заметки")
			}
			if weekDayInt < 1 || weekDayInt > 7 {
				return "", fmt.Errorf("недопустимый диапазон дней недели для переноса заметки")
			}
			weekDaysInt = append(weekDaysInt, weekDayInt)
		}

		nextDate = now

		for _, weekDay := range weekDaysInt {
			if now.Weekday() <= time.Weekday(weekDay) {
				nextDate = nextDate.AddDate(0, 0, int(time.Weekday(weekDay)-now.Weekday()))
				break
			}
		}

		if nextDate == now {
			nextDate = nextDate.AddDate(0, 0, 7-int(now.Weekday())+weekDaysInt[0])
		}
	}

	return nextDate.Format(DateFormat), nil
}
