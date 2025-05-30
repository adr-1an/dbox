package dbox

import (
	"database/sql"
	"fmt"
	"os"
)

func Rollback(dbType, dsn string, pretend bool) {
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		fmt.Println("Failed to connect to DB:", err)
		os.Exit(1)
	}
	defer db.Close()

	var latest string
	err = db.QueryRow("SELECT name FROM migrations ORDER BY name DESC LIMIT 1").Scan(&latest)
	if err == sql.ErrNoRows {
		fmt.Println("No migrations to rollback.")
		return
	}
	if err != nil {
		fmt.Println("Failed to fetch latest migration:", err)
		fmt.Println("Try running ./dbox init")
		os.Exit(1)
	}

	downPath := "db/migrations/" + latest + "/down.sql"
	downSQL, err := os.ReadFile(downPath)
	if err != nil {
		fmt.Println("Failed to read down.sql for", latest)
		os.Exit(1)
	}

	if pretend {
		fmt.Printf("== Pretending rollback: %s\n-- SQL:\n%s\n\n", latest, string(downSQL))
		return
	}

	_, err = db.Exec(string(downSQL))
	if err != nil {
		fmt.Println("Rollback failed:", latest)
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	del := fmt.Sprintf("DELETE FROM migrations WHERE name = %s", placeholder(dbType, 1))
	_, err = db.Exec(del, latest)
	if err != nil {
		fmt.Println("Failed to remove migration record:", latest)
		os.Exit(1)
	}

	fmt.Println("Rolled back:", latest)
}
