package dbox

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"os"
	"strings"
)

func Init(dbType, dsn string) {
	if dbType != "sqlite" {
		ensureEnvVar("DB_USER")
		ensureEnvVar("DB_PASS")
		ensureEnvVar("DB_HOST")
		ensureEnvVar("DB_PORT")
		ensureEnvVar("DB_NAME")
	} else {
		ensureEnvVar("DB_DATABASE")
	}

	// Create db/ dir
	err := os.MkdirAll("db", 0755)
	if err != nil {
		fmt.Println("Failed to create db directory:", err)
		os.Exit(1)
	}

	// Create db/migrations/ dir
	err = os.MkdirAll("db/migrations", 0755)
	if err != nil {
		fmt.Println("Failed to create migration directory:", err)
	}

	// Ask user if they want to create migrations table in the db
	var input string
	fmt.Println("Create migration records table? (Table name will be \"migrations\")")
	fmt.Println("(You can skip this if the table already exists)")
	fmt.Print("[yes/no] > ")
	fmt.Scanln(&input)
	input = strings.ToLower(input)
	if input == "yes" {
		db, err := sql.Open(dbType, dsn)
		if err != nil {
			fmt.Println("Failed to connect to db:", err)
			os.Exit(1)
		}
		defer db.Close()

		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS migrations (
			    name VARCHAR(255) PRIMARY KEY
			)
		`)
		if err != nil {
			fmt.Println("Failed to create migration table:", err)
			os.Exit(1)
		}
		fmt.Println("Migration records table created.")
	}

	fmt.Println("Initialized DBox project.")
}

func ensureEnvVar(key string) {
	_, exists := os.LookupEnv(key)
	if exists {
		return
	}

	var input string
	fmt.Printf("Missing %s in .env\nWould you like to configure it now?\n", key)
	fmt.Print("[yes/no] > ")
	fmt.Scanln(&input)
	input = strings.ToLower(input)

	if input != "yes" {
		fmt.Printf("Canceled. %s is required.\n", key)
		os.Exit(1)
	}

	fmt.Printf("[.env] %s: > ", key)
	var value string
	fmt.Scanln(&value)

	appendToEnv(key, value)
}

func appendToEnv(key, value string) {
	file, err := os.OpenFile(".env", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Failed to open .env:", err)
		os.Exit(1)
	}
	defer file.Close()

	line := fmt.Sprintf("%s=%s\n", key, value)
	file.WriteString("\n" + line)
}
