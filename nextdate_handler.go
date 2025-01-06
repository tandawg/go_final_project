package main

import (
	"fmt"
	"net/http"
	"time"
	"strings"
)

// Функция обработчика для маршрута /api/nextdate
// Вычисляет следующую дату на основе заданного правила повторения
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Преобразуем строку now в time.Time
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Некорректный формат даты 'now'", http.StatusBadRequest)
		return
	}

	// Пробуем преобразовать date в time.Time
	dateParsed, err := time.Parse("20060102", date)
	if err != nil || !isValidDate(date) { // Используем isValidDate
		// Если date некорректна или невалидна, возвращаем пустую строку
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "")
		return
	}

	// Обрабатываем правило repeat
	switch {
	case repeat == "y":
		// Добавляем годы, пока не достигнем ближайшей даты после now
		nextDate := dateParsed.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, nextDate.Format("20060102"))

	case strings.HasPrefix(repeat, "d "):
		var days int
		_, err := fmt.Sscanf(repeat, "d %d", &days)
		if err != nil || days <= 0 || days > 400 {
			// Некорректное количество дней
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprint(w, "")
			return
		}
		nextDate := dateParsed.AddDate(0, 0, days)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, nextDate.Format("20060102"))

	default:
		// Если repeat пустой или содержит некорректное значение
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "")
	}
}