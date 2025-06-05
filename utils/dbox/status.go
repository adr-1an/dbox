package dbox

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gocql/gocql"
)

func Status(dbType, dsn string) {
	fmt.Println("Migration Status\n")

	applied := make(map[string]bool)

	if dbType == "cql" {
		session, err := openCQLSession()
		if err != nil {
			fmt.Println("Error opening database:", err)
			os.Exit(1)
		}

		iter := session.Query("SELECT name FROM migrations").Iter()
		var n string
		for iter.Scan(&n) {
			applied[n] = true
		}
		if err := iter.Close(); err != nil {
			fmt.Println("Failed to fetch applied migrations:", err)
			session.Close()
			os.Exit(1)
		}
		session.Close()
	} else {
		db, err := sql.Open(dbType, dsn)
		if err != nil {
			fmt.Println("Error opening database:", err)
			os.Exit(1)
		}
		defer db.Close()

		rows, err := db.Query("SELECT name FROM migrations")
		if err != nil {
			fmt.Println("Failed to fetch applied migrations:", err)
			os.Exit(1)
		}
		defer rows.Close()

		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				fmt.Println("Scan error:", err)
				os.Exit(1)
			}
			applied[name] = true
		}
	}

	entries, err := os.ReadDir("db/migrations")
	if len(entries) == 0 {
		fmt.Println("No migrations found in db/migrations.")
		return
	}
	if err != nil {
		fmt.Println("Failed to read migrations dir:", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		if applied[name] {
			fmt.Printf("\033[32m[✓]\033[0m %s\n", name)
		} else {
			fmt.Printf("\033[31m[✗]\033[0m %s\n", name)
		}
	}
}
