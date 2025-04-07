package dbox

import (
	"fmt"
	"os"
	"time"
)

func Create() {
	if len(os.Args) < 3 {
		fmt.Println("You need to give the migration a name, like: ./dbox create create_users_table")
		os.Exit(0)
	}

	name := os.Args[2]
	timestamp := time.Now().Format("20060102_1504")
	dirName := fmt.Sprintf("db/migrations/%s_%s", timestamp, name)

	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		os.Exit(1)
	}

	upFile := dirName + "/up.sql"
	downFile := dirName + "/down.sql"

	os.WriteFile(upFile, []byte("-- Up"), 0644)
	os.WriteFile(downFile, []byte("-- Down"), 0644)

	fmt.Println("Created migration at:", dirName)
}
