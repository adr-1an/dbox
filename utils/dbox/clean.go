package dbox

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func Clean(dbType, dsn string) {
	var input string
	fmt.Println("This will delete all migration records that no longer have folders.")
	fmt.Println("Are you sure you want to continue?")
	fmt.Print("[yes/no] > ")
	fmt.Scanln(&input)
	input = strings.ToLower(input)
	if input != "yes" {
		fmt.Println("Canceled.")
		os.Exit(0)
	}

	var deleted int
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		fmt.Println("Failed to connect to DB:", err)
		os.Exit(1)
	}
	defer db.Close()

	// STEP 1: fetch all migration names into a slice
	rows, err := db.Query("SELECT name FROM migrations")
	if err != nil {
		fmt.Println("Failed to query migrations:", err)
		os.Exit(1)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			fmt.Println("Scan error:", err)
			os.Exit(1)
		}
		names = append(names, name)
	}
	rows.Close()

	for _, name := range names {
		path := "db/migrations/" + name

		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			_, err := db.Exec("DELETE FROM migrations WHERE name = ?", name)
			if err != nil {
				fmt.Println("Failed to delete migration record:", err)
				os.Exit(1)
			}

			fmt.Println("Deleted:", name)
			deleted++
		}
	}

	if deleted > 0 {
		word := "records"
		if deleted == 1 {
			word = "record"
		}
		fmt.Println("\nSuccess!\n")
		fmt.Printf("Deleted %d migration %s.\n", deleted, word)
	} else {
		fmt.Println("\nThere are no migration records to delete. Nothing was changed.")
	}
}
