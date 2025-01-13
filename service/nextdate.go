package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// const (
// 	DateFormat      = "20060102"
// 	transferDaysMax = 400
// )

var DateFormat = "20060102"
var TransferDaysMax = 400

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" || repeat == " " {
		// err := fmt.Errorf("правило повторения не должно содержать пустую строку")
		// return "", err
		return "", fmt.Errorf("правило повторения не должно содержать пустую строку")
	}

	nextDate, err := time.Parse(DateFormat, date)
	if err != nil {
		//return "", fmt.Errorf("не верный формат даты: %s", date)
		return "", err
	}

	// Определеяем правило повторения задач
	transferRule := strings.Split(repeat, "")[0]
	if !(strings.EqualFold(transferRule, "d") || strings.EqualFold(transferRule, "y")) {
		return "", fmt.Errorf("неподдерживаемый формат")
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

	}

	return nextDate.Format(DateFormat), nil
}

// switch strings.ToLower(repeatRule) {
// case "d", "y", "w", "m":

// 	switch strings.ToLower(repeatRule) {
// 	// обарботка правила повторения тип "d" - день
// 	case "d":
// 		splitedRepeatRule := strings.Split(repeat, " ")
// 		daysTransfer, err := strconv.Atoi(splitedRepeatRule[1])
// 		if err != nil {
// 			err := fmt.Errorf("не верный формат количества дней переноса заметки: %s", splitedRepeatRule[1])
// 			return "", err
// 		}

// 		if daysTransfer < 1 || daysTransfer > TransferDaysMax {
// 			err := fmt.Errorf("количество дней для переноса должно быть от 1 до %d", transferDaysMax)
// 			return "", err
// 		}

// 		nextDate = nextDate.AddDate(0, 0, daysTransfer)
// 		for nextDate.Before(now) {
// 			nextDate = nextDate.AddDate(0, 0, daysTransfer)
// 		}

// 	// обарботка правила повторения тип "y" - год
// 	case "y":
// 		nextDate = nextDate.AddDate(1, 0, 0)
// 		for nextDate.Before(now) {
// 			nextDate = nextDate.AddDate(1, 0, 0)
// 		}

// 	// обарботка правила повторения тип "w" - дни недели
// 	case "w":
// 		err := fmt.Errorf("не поддерживаемый формат правила повторения задачи")
// 		return "", err

// 	// обарботка правила повторения тип "m" - месяцы
// 	case "m":
// 		err := fmt.Errorf("не поддерживаемый формат правила повторения задачи")
// 		return "", err
// 	}

// default:
// 	return "", fmt.Errorf("некорректный формат параметра repeat: %s", repeat)
// }

// return nextDate.Format(DateFormat), nil
