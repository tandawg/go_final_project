package main

import (
	"encoding/json"
	"net/http"
	"time"
	"database/sql"
	"log"
)

// Обработчик для маршрута /api/tasks
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {

	// проверяем, что метод запроса - GET
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	// заголовок ответа
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Чтение параметра поиска из запроса
search := r.URL.Query().Get("search")
var tasks []Task

// Формирование SQL-запроса
var rows *sql.Rows
var err error

if search != "" {
    // Попытка распарсить параметр как дату
    date, err := time.Parse("02.01.2006", search)
    if err == nil { // Если это дата
        // Преобразуем дату в нужный формат 20060102
        search = date.Format("20060102")
        // Ищем задачи по дате
        rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT 50`, search)
        if err != nil {
            log.Fatal(err) // Ошибка при запросе
        }
    } else { // Если это строка для поиска в заголовке или комментарии
        rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT 50`, "%"+search+"%", "%"+search+"%")
        if err != nil {
            log.Fatal(err) // Ошибка при запросе
        }
    }
} else {
    // Если search не указан, просто берем ближайшие задачи
    rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT 50`)
    if err != nil {
        log.Fatal(err) // Ошибка при запросе
    }
}

// Обработка результатов запроса
defer rows.Close()
for rows.Next() {
    var task Task
    if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка чтения данных"})
        return
    }
    tasks = append(tasks, task)
}

// Если задач нет, возвращаем пустой список
if len(tasks) == 0 {
    tasks = []Task{}
}

// Формируем и отправляем ответ
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string][]Task{"tasks": tasks})
}