package dbox

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func Refresh(dbType, dsn string, pretend bool) {
	var input string
	fmt.Println("WARNING: This will clear the entire database.")
	fmt.Println("Are you sure you want to continue?")
	fmt.Print("[yes/no] > ")
	fmt.Scanln(&input)
	input = strings.ToLower(input)
	if input != "yes" {
		fmt.Println("Canceled.")
		os.Exit(0)
	}

	db, err := sql.Open(dbType, dsn)
	if err != nil {
		fmt.Println("Failed to connect to DB:", err)
		os.Exit(1)
	}

	rows, err := db.Query("SELECT name FROM migrations ORDER BY name DESC")
	if err != nil {
		fmt.Println("Failed to fetch migrations:", err)
		os.Exit(1)
	}

	var names []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			fmt.Println("Failed to scan migration name:", err)
			os.Exit(1)
		}
		names = append(names, name)
	}

	rows.Close()
	db.Close()

	for _, name := range names {
		db, err := sql.Open("sqlite", "db/database.sqlite")
		if err != nil {
			fmt.Println("Failed to connect to DB:", err)
			os.Exit(1)
		}

		downPath := "db/migrations/" + name + "/down.sql"
		sqlBytes, err := os.ReadFile(downPath)
		if err != nil {
			fmt.Println("Failed to read:", downPath)
			os.Exit(1)
		}

		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			fmt.Println("Failed to rollback:", name)
			fmt.Println("Error:", err)
			db.Close()
			os.Exit(1)
		}

		_, err = db.Exec("DELETE FROM migrations WHERE name = ?", name)
		if err != nil {
			fmt.Println("Failed to delete migration record:", name)
			db.Close()
			os.Exit(1)
		}

		db.Close()
		fmt.Println("Rolled back:", name)
	}
	Migrate(dbType, dsn, pretend)
}
