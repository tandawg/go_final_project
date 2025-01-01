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
		if nextDate.Before(now) {
			for nextDate.Before(now) {
				nextDate = nextDate.AddDate(0, 0, days)
				if nextDate.Year() > 2100 {
					return "", nil
				}
			}
		}
		return nextDate.Format("20060102"), nil

	case "y": // Повторение в годах
		nextDate := taskDate.AddDate(1, 0, 0)
		if nextDate.Before(now) {
			for nextDate.Before(now) {
				nextDate = nextDate.AddDate(1, 0, 0)
				if nextDate.Year() > 2100 {
					return "", nil
				}
			}
		}
		return nextDate.Format("20060102"), nil

	default:
		return "", nil
	}
}