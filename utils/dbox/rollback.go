package dbox

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
)

func Rollback(dbType, dsn string, pretend bool) {
	if dbType == "cql" {
		session, err := openCQLSession()
		if err != nil {
			fmt.Println("Failed to connect to DB:", err)
			os.Exit(1)
		}
		defer session.Close()

		var names []string
		iter := session.Query("SELECT name FROM migrations").Iter()
		var n string
		for iter.Scan(&n) {
			names = append(names, n)
		}
		if err := iter.Close(); err != nil {
			fmt.Println("DB error:", err)
			os.Exit(1)
		}
		if len(names) == 0 {
			fmt.Println("No migrations to rollback.")
			return
		}
		sort.Strings(names)
		latest := names[len(names)-1]

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

		if err := session.Query(string(downSQL)).Exec(); err != nil {
			fmt.Println("Rollback failed:", latest)
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		if err := session.Query("DELETE FROM migrations WHERE name = ?", latest).Exec(); err != nil {
			fmt.Println("Failed to remove migration record:", latest)
			os.Exit(1)
		}

		fmt.Println("Rolled back:", latest)
		return
	}

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
