package main

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func createDatabase() *sql.DB {
	dbFile := "scheduler.db"
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil && os.IsNotExist(err) {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if install {
		createTableQuery := `
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT CHECK(LENGTH(repeat) <= 128)
		);
		CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
		`
		_, err = db.Exec(createTableQuery)
		if err != nil {
			log.Fatalf("Ошибка создания схемы базы данных: %v", err)
		}
		log.Println("База данных и таблица успешно созданы.")
	}
	return db
}