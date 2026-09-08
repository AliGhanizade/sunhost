package config

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./sunhost.db")
	if err != nil {
		log.Fatal(err)
	}
	// File-backed SQLite allows only one writer at a time; limiting the pool
	// to a single connection avoids "database is locked" errors.
	DB.SetMaxOpenConns(1)
	if _, err := DB.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		log.Fatal(err)
	}
}

func MigrateUser() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			full_name TEXT NOT NULL,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func MigrateSystemStat() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS system_stats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			alloc_ram REAL NOT NULL,
			goroutines INTEGER NOT NULL,
			live_objects INTEGER NOT NULL,
			record_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal("error creating system_stats table: ", err)
	}
}

func MigrateUserLog() {
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip_address TEXT NOT NULL,
		system_info TEXT NOT NULL,
		action TEXT NOT NULL,
		time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		username TEXT NOT NULL,
		FOREIGN KEY (username) REFERENCES users(username)
	);`)

	if err != nil {
		log.Fatal(err)
	}
}
