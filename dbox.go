package main

import (
	"app/utils/dbox"
	"fmt"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file not found or failed to load.")
	}

	if os.Getenv("DB_TYPE") == "" {
		fmt.Println("You need to set DB_TYPE in your .env file (sqlite, mysql, or postgres)")
		os.Exit(1)
	}

	pretend := false
	for _, arg := range os.Args {
		if arg == "--pretend" || arg == "--dry-run" || arg == "-p" {
			pretend = true
			break
		}
	}

	// Get settings from .env
	godotenv.Load()
	dbType := os.Getenv("DB_TYPE")
	dsn := BuildDSN(dbType)

	if len(os.Args) < 2 {
		fmt.Println("No command provided. Type ./dbox help to view all commands.")
		os.Exit(0)
	}

	cmd := os.Args[1]

	switch cmd {
	case "create", "make":
		dbox.Create()

	case "migrate", "up":
		dbox.Migrate(dbType, dsn, pretend)

	case "rollback", "down":
		dbox.Rollback(dbType, dsn, pretend)

	case "refresh":
		dbox.Refresh(dbType, dsn, pretend)

	case "status", "stats", "stat":
		dbox.Status(dbType, dsn)

	case "clean":
		dbox.Clean(dbType, dsn)

	case "init", "initialize":
		dbox.Init(dbType, dsn)

	case "help", "?":
		fmt.Println("DBox - DB Toolbox\n")
		fmt.Println("./dbox help - Shows this menu.")
		fmt.Println("./dbox create [migration_name] - Creates a migration with the specified name.")
		fmt.Println("./dbox migrate - Run all migrations.\n    --pretend (-p for short) - shows the SQL that would run, but doesn't execute it\n       --dry-run achieves the same result.")
		fmt.Println("./dbox rollback - Roll the last migration back.")
		fmt.Println("./dbox refresh - Re-run all migrations' up and down methods.")
		fmt.Println("./dbox clean - Deletes all migration registries that don't have a corresponding directory.")
		fmt.Println("./dbox init - Initializes the main file structure:\n    - Creates db/ and db/migrations")
	default:
		fmt.Println("Unknown command:", cmd)
		fmt.Println("Run ./dbox help for more info.")
	}
}

func BuildDSN(dbType string) string {
	switch dbType {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)

	case "postgres":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_NAME"),
		)

	case "sqlite":
		return os.Getenv("DB_DATABASE")

	default:
		fmt.Println("Unsupported DB_TYPE:", dbType)
		os.Exit(1)
		return ""
	}
}
