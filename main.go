package main

func main() {
	// Запуск базы данных
	db := createDatabase()
	defer db.Close()

	// Запуск веб-сервера
	startServer()
}