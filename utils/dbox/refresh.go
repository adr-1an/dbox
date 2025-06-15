package dbox

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
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

	if dbType == "cql" {
		session, cleanup, err := cqlConn()
		if err != nil {
			fmt.Println("Failed to connect to DB:", err)
			os.Exit(1)
		}
		defer cleanup()

		iter := session.Query("SELECT name FROM migrations").Iter()
		var names []string
		var n string
		for iter.Scan(&n) {
			names = append(names, n)
		}
		if err := iter.Close(); err != nil {
			fmt.Println("Failed to fetch migrations:", err)
			os.Exit(1)
		}

		sort.Sort(sort.Reverse(sort.StringSlice(names)))

		for _, name := range names {
			sess, cqlCleanup, err := cqlConn()
			if err != nil {
				fmt.Println("Failed to connect to DB:", err)
				os.Exit(1)
			}
			downPath := "db/migrations/" + name + "/down.sql"
			sqlBytes, err := os.ReadFile(downPath)
			if err != nil {
				fmt.Println("Failed to read:", downPath)
				cqlCleanup()
				os.Exit(1)
			}
			if err := sess.Query(string(sqlBytes)).Exec(); err != nil {
				fmt.Println("Failed to rollback:", name)
				fmt.Println("Error:", err)
				cqlCleanup()
				os.Exit(1)
			}
			if err := sess.Query("DELETE FROM migrations WHERE name = ?", name).Exec(); err != nil {
				fmt.Println("Failed to delete migration record:", name)
				cqlCleanup()
				os.Exit(1)
			}
			cqlCleanup()
			fmt.Println("Rolled back:", name)
		}

		Migrate(dbType, dsn, pretend)
		return
	}

	db, cleanup, err := sqlConn(dbType, dsn)
	if err != nil {
		fmt.Println("Failed to connect to DB:", err)
		os.Exit(1)
	}
	defer cleanup()

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

	for _, name := range names {
		db, c, err := sqlConn(dbType, dsn)
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
			c()
			os.Exit(1)
		}

		del := fmt.Sprintf("DELETE FROM migrations WHERE name = %s", placeholder(dbType, 1))
		_, err = db.Exec(del, name)
		if err != nil {
			fmt.Println("Failed to delete migration record:", name)
			c()
			os.Exit(1)
		}

		c()
		fmt.Println("Rolled back:", name)
	}
	Migrate(dbType, dsn, pretend)
}
