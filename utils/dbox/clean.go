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
	if dbType == "cql" {
		session, err := openCQLSession()
		if err != nil {
			fmt.Println("Failed to connect to DB:", err)
			os.Exit(1)
		}

		iter := session.Query("SELECT name FROM migrations").Iter()
		var names []string
		var n string
		for iter.Scan(&n) {
			names = append(names, n)
		}
		if err := iter.Close(); err != nil {
			fmt.Println("Failed to query migrations:", err)
			session.Close()
			os.Exit(1)
		}

		for _, name := range names {
			path := "db/migrations/" + name
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				if err := session.Query("DELETE FROM migrations WHERE name = ?", name).Exec(); err != nil {
					fmt.Println("Failed to delete migration record:", err)
					session.Close()
					os.Exit(1)
				}
				fmt.Println("Deleted:", name)
				deleted++
			}
		}
		session.Close()

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
		return
	}

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
			del := fmt.Sprintf("DELETE FROM migrations WHERE name = %s", placeholder(dbType, 1))
			_, err := db.Exec(del, name)
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
