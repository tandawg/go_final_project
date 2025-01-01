package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Функция обработчика для маршрута /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Проверка на пустое значение для repeat
	if repeat == "" {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "") // Если repeat пуст, возвращаем пустую строку
		return
	}

	// Преобразуем строку now в time.Time
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при разборе даты 'now': %v", err), http.StatusBadRequest)
		return
	}

	// Используем переменную now для логирования
	log.Printf("Текущая дата (now): %v\n", now)

	// Преобразуем строку date в time.Time
	dateParsed, err := time.Parse("20060102", date)
	if err != nil {
		http.Error(w, fmt.Sprintf("Некорректный формат даты: %v", date), http.StatusBadRequest)
		return
	}

	// Если repeat == "y", сдвигаем дату на 1 год вперед
	if repeat == "y" {
		nextDate := dateParsed.AddDate(1, 0, 0)
		if nextDate.Before(dateParsed) {
			http.Error(w, fmt.Sprintf("Следующая дата %v не может быть раньше текущей даты %v", nextDate.Format("20060102"), dateParsed.Format("20060102")), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, nextDate.Format("20060102"))
		return
	}

	// Если repeat начинается с "d", сдвигаем дату на определенное количество дней
	if len(repeat) > 1 && repeat[:1] == "d" {
		// Парсим количество дней
		var days int
		_, err := fmt.Sscanf(repeat, "d %d", &days)
		if err != nil {
			http.Error(w, fmt.Sprintf("Некорректное значение дней: %v", repeat), http.StatusBadRequest)
			return
		}
		nextDate := dateParsed.AddDate(0, 0, days)
		if nextDate.Before(dateParsed) {
			http.Error(w, fmt.Sprintf("Следующая дата %v не может быть раньше текущей даты %v", nextDate.Format("20060102"), dateParsed.Format("20060102")), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, nextDate.Format("20060102"))
		return
	}

	// Если repeat содержит другие значения, возвращаем ошибку
	http.Error(w, fmt.Sprintf("Неподдерживаемое правило повторения: %v", repeat), http.StatusBadRequest)
}

func startServer() {
	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	// Логирование для проверки пути
	staticPath := "./go_final_project/web"
	log.Println("Serving static files from:", staticPath)

	// Обработчик для других статических файлов
	http.Handle("/", http.FileServer(http.Dir(staticPath)))

	// Обработчик для API /api/nextdate
	http.HandleFunc("/api/nextdate", nextDateHandler)

	fmt.Printf("Сервер работает на порту %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}
}