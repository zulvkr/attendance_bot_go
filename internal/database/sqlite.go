package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// SQLiteDB wraps the sql.DB connection
type SQLiteDB struct {
	*sql.DB
}

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	sqliteDB := &SQLiteDB{DB: db}

	// Initialize schema
	if err := sqliteDB.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return sqliteDB, nil
}

// initSchema creates the necessary tables and indexes
func (db *SQLiteDB) initSchema() error {
	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Create attendance table
	attendanceTableSQL := `
	CREATE TABLE IF NOT EXISTS attendance (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		username TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT,
		timestamp TEXT NOT NULL,
		type TEXT NOT NULL CHECK (type IN ('check_in', 'check_out')),
		date TEXT NOT NULL,
		UNIQUE(user_id, date, type)
	);`

	if _, err := db.Exec(attendanceTableSQL); err != nil {
		return fmt.Errorf("failed to create attendance table: %w", err)
	}

	// Create indexes for attendance table
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_user_date ON attendance(user_id, date);",
		"CREATE INDEX IF NOT EXISTS idx_date ON attendance(date);",
		"CREATE INDEX IF NOT EXISTS idx_user_id ON attendance(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_type ON attendance(type);",
	}

	for _, indexSQL := range indexes {
		if _, err := db.Exec(indexSQL); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	// Create alias table
	aliasTableSQL := `
	CREATE TABLE IF NOT EXISTS alias (
		user_id INTEGER PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT
	);`

	if _, err := db.Exec(aliasTableSQL); err != nil {
		return fmt.Errorf("failed to create alias table: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *SQLiteDB) Close() error {
	return db.DB.Close()
}

// BeginTx starts a new transaction
func (db *SQLiteDB) BeginTx() (*sql.Tx, error) {
	return db.DB.Begin()
}
