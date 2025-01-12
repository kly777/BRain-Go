package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./brain.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		password TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS cards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS relations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		relata INTEGER NOT NULL,
		relation INTEGER NOT NULL,
		position INTEGER NOT NULL,
		describe TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY(relata) REFERENCES cards(id),
		FOREIGN KEY(relation) REFERENCES cards(id)
	);`

	_, err = DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}
