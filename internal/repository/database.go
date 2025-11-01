package repository

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
)

// Database представляет подключение к базе данных
type Database struct {
	conn *sql.DB
	psql squirrel.StatementBuilderType
}

// NewDatabase создает новое подключение к базе данных
func NewDatabase(dbPath string) (*Database, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &Database{
		conn: conn,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}

	if err := db.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// createTables создает необходимые таблицы
func (db *Database) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			username TEXT,
			first_name TEXT NOT NULL,
			last_name TEXT,
			chat_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(id, chat_id)
		)`,
		`CREATE TABLE IF NOT EXISTS person_of_the_day (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			chat_id INTEGER NOT NULL,
			date DATE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			UNIQUE(chat_id, date)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_person_of_the_day_chat_date ON person_of_the_day(chat_id, date)`,
		`CREATE INDEX IF NOT EXISTS idx_users_chat ON users(chat_id)`,
	}

	for _, query := range queries {
		if _, err := db.conn.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %s: %w", query, err)
		}
	}

	return nil
}

// Close закрывает соединение с базой данных
func (db *Database) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}
