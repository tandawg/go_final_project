package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate рассчитывает следующую дату задачи с учетом правила повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Проверка длины строки даты
	if len(date) != 8 {
		fmt.Printf("Отказ: некорректная дата %s\n", date)
		return "", nil
	}

	// Извлекаем год, месяц и день из строки
	year := date[:4]
	month := date[4:6]
	day := date[6:8]

	// Парсим год, месяц и день
	yearNum, err := strconv.Atoi(year)
	if err != nil {
		fmt.Printf("Отказ: некорректный год %s\n", year)
		return "", nil
	}

	monthNum, err := strconv.Atoi(month)
	if err != nil || monthNum < 1 || monthNum > 12 {
		fmt.Printf("Отказ: некорректный месяц %s\n", month)
		return "", nil
	}

	dayNum, err := strconv.Atoi(day)
	if err != nil || dayNum < 1 || dayNum > 31 {
		fmt.Printf("Отказ: некорректный день %s\n", day)
		return "", nil
	}

	// Создаем исходную дату
	taskDate := time.Date(yearNum, time.Month(monthNum), dayNum, 0, 0, 0, 0, time.UTC)

	// Если год некорректен (меньше 1900), заменяем на ближайший корректный
	if taskDate.Year() < 1900 {
		taskDate = time.Date(2024, taskDate.Month(), taskDate.Day(), 0, 0, 0, 0, time.UTC)
		fmt.Printf("Год %d был заменен на ближайший корректный: %d\n", yearNum, taskDate.Year())
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
		if err != nil || days <= 0 || days > 400 {
			fmt.Printf("Отказ: некорректное количество дней %s\n", repeatParts[1])
			return "", nil
		}

		nextDate := taskDate.AddDate(0, 0, days)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}

		if nextDate.Year() > 2100 {
			return "", nil
		}
		return nextDate.Format("20060102"), nil

	case "y": // Повторение в годах
		// Прибавляем год, пока дата не станет больше текущей
		nextDate := taskDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}

		if nextDate.Year() > 2100 {
			return "", nil
		}
		return nextDate.Format("20060102"), nil

	default:
		fmt.Printf("Отказ: некорректное правило повторения %s\n", repeat)
		return "", nil
	}
}