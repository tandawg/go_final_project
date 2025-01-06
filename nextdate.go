package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate рассчитывает следующую дату задачи с учетом правила повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Парсинг даты задачи
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("некорректный формат даты: %s", date)
	}

	// Проверка диапазона допустимых годов
	if taskDate.Year() < 1900 || taskDate.Year() > 2100 {
		fmt.Printf("Отказ: год %d вне допустимого диапазона (1900–2100)\n", taskDate.Year())
		return "", nil // Возвращаем пустую строку для некорректных годов
	}

	// Разбиение правила повторения
	repeatParts := strings.Fields(repeat)
	if len(repeatParts) == 0 {
		return "", nil
	}

	switch repeatParts[0] {
	case "d": // Повторение в днях
		if len(repeatParts) != 2 {
			return "", nil
		}
		days, err := strconv.Atoi(repeatParts[1])
		if err != nil || days <= 0 {
			return "", nil
		}

		nextDate := taskDate.AddDate(0, 0, days)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, days)
			if nextDate.Year() > 2100 {
				return "", nil
			}
		}
		return nextDate.Format("20060102"), nil

	case "y": // Повторение в годах
		nextDate := taskDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
			if nextDate.Year() > 2100 {
				return "", nil
			}
		}
		return nextDate.Format("20060102"), nil

	default:
		return "", nil
	}
}

// Вспомогательная функция, проверяет корректность формата даты в строке (YYYYMMDD)
func isValidDate(date string) bool {
	if len(date) != 8 {
		return false
	}
	year := date[:4]
	month := date[4:6]
	day := date[6:8]

	// Проверка все ли части даты валидны
	yearNum, err1 := strconv.Atoi(year)
	monthNum, err2 := strconv.Atoi(month)
	dayNum, err3 := strconv.Atoi(day)

	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}

	// Проверка по верхней границе дат, включая исторические даты
	if yearNum > 2100 || monthNum < 1 || monthNum > 12 || dayNum < 1 || dayNum > 31 {
		return false
	}

	return true
}