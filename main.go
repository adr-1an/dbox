package main

import (
        "bufio"
        "database/sql"
        "dbox/utils/dbox"
        "fmt"
        "os"
        "strings"

        "github.com/joho/godotenv"
        _ "modernc.org/sqlite"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file not found or failed to load.")
	}

	if os.Getenv("DB_TYPE") == "" {
		fmt.Println("You need to set DB_TYPE in your .env file (sqlite, mysql, postgres, clickhouse, or cql)")
		os.Exit(1)
	}

        // Get settings from .env
        godotenv.Load()
        dbType := os.Getenv("DB_TYPE")
        dsn := BuildDSN(dbType)

	if os.Getenv("DBOX_TYPE") == "console" {
		runConsole(dbType, dsn)
		return
	}

	runCommand(os.Args[1:], dbType, dsn)
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

	case "clickhouse":
		return fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)

	case "cql":
		return fmt.Sprintf("%s:%s/%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)

	default:
		fmt.Println("Unsupported DB_TYPE:", dbType)
		os.Exit(1)
		return ""
	}
}

func runConsole(dbType, dsn string) {
	fmt.Println("DBox Console - type 'exit' to quit")
	if dbType == "cql" {
		sess, err := dbox.OpenCQLSession()
		if err != nil {
			fmt.Println("Failed to connect to DB:", err)
			os.Exit(1)
		}
		dbox.SharedSession = sess
		defer func() {
			sess.Close()
			dbox.SharedSession = nil
		}()
	} else {
		db, err := sql.Open(dbType, dsn)
		if err != nil {
			fmt.Println("Failed to connect to DB:", err)
			os.Exit(1)
		}
		dbox.SharedDB = db
		defer func() {
			db.Close()
			dbox.SharedDB = nil
		}()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("dbox> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" {
			break
		}
		args := strings.Fields(line)
		runCommand(args, dbType, dsn)
	}
}

func runCommand(args []string, dbType, dsn string) {
	if len(args) < 1 {
		fmt.Println("No command provided. Type help to view all commands.")
		return
	}

	pretend := false
	for _, arg := range args {
		if arg == "--pretend" || arg == "--dry-run" || arg == "-p" {
			pretend = true
			break
		}
	}

	os.Args = append([]string{os.Args[0]}, args...)

	cmd := args[0]

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
		fmt.Println("DBox - DB Toolbox")
		fmt.Println("./dbox help - Shows this menu.")
		fmt.Println("./dbox create [migration_name] - Creates a migration with the specified name.")
		fmt.Println("./dbox migrate - Run all migrations.\n    --pretend (-p for short) - shows the SQL that would run, but doesn't execute it\n    --dry-run achieves the same result.")
		fmt.Println("./dbox rollback - Roll the last migration back.")
		fmt.Println("./dbox refresh - Re-run all migrations from scratch.")
		fmt.Println("./dbox clean - Deletes all migration registries that don't have a corresponding directory.")
		fmt.Println("./dbox init - Initializes the main file structure:\n    - Creates db/ and db/migrations")
	default:
		fmt.Println("Unknown command:", cmd)
		fmt.Println("Run ./dbox help for more info.")
	}
}
