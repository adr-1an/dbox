package dbox

import (
	"database/sql"
	"fmt"
	"os"
)

func Migrate(dbType, dsn string, pretend bool) {
	// Create the db dir if it doesn't exist
	if _, err := os.Stat("db"); os.IsNotExist(err) {
		os.Mkdir("db", 0755)
	}
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		fmt.Println("Failed to connect to DB:", err)
	}
	defer db.Close()

	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS migrations (
			    name TEXT PRIMARY KEY
			)
		`)
	if err != nil {
		fmt.Println("Failed to create migrations table:", err)
		os.Exit(1)
	}

	entries, err := os.ReadDir("db/migrations")
	if err != nil {
		fmt.Println("Failed to read migrations folder:", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		migrationName := entry.Name()

		var exists string
		query := fmt.Sprintf("SELECT name FROM migrations WHERE name = %s", placeholder(dbType, 1))
		err := db.QueryRow(query, migrationName).Scan(&exists)
		if err != nil && err != sql.ErrNoRows {
			fmt.Println("DB error:", err)
			os.Exit(1)
		}

		if exists == migrationName {
			continue
		}

		upPath := "db/migrations/" + migrationName + "/up.sql"
		sqlBytes, err := os.ReadFile(upPath)
		if err != nil {
			fmt.Println("Failed to read:", upPath)
			os.Exit(1)
		}

		if pretend {
			fmt.Printf("== Pretending migration: %s\n-- SQL:\n%s\n\n", migrationName, string(sqlBytes))
			continue
		} else {
			_, err = db.Exec(string(sqlBytes))
			if err != nil {
				fmt.Println("Migration failed:", migrationName)
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}

		insert := fmt.Sprintf("INSERT INTO migrations (name) VALUES (%s)", placeholder(dbType, 1))
		_, err = db.Exec(insert, migrationName)
		if err != nil {
			fmt.Println("Failed to record migration:", migrationName)
			os.Exit(1)
		}

		fmt.Println("Migrated:", migrationName)
	}
}
