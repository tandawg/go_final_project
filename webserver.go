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
	_, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при разборе даты 'now': %v", err), http.StatusBadRequest)
		return
	}

	// Преобразуем строку date в time.Time
	dateParsed, err := time.Parse("20060102", date)
	if err != nil {
		http.Error(w, fmt.Sprintf("Некорректный формат даты: %v", date), http.StatusBadRequest)
		return
	}

	// Проверка диапазона годов
	if dateParsed.Year() < 1900 || dateParsed.Year() > 2100 {
		log.Printf("Год %d выходит за пределы диапазона (1900–2100)", dateParsed.Year())
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "") // Возвращаем пустую строку для недопустимых годов
		return
	}

	// Если repeat == "y", сдвигаем дату на 1 год вперед
	if repeat == "y" {
		nextDate := dateParsed.AddDate(1, 0, 0)
		log.Printf("Повторение год: %v -> %v\n", dateParsed, nextDate)
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
		log.Printf("Повторение дни: %v + %d дней -> %v\n", dateParsed, days, nextDate)
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