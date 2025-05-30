package dbox

import "fmt"

func placeholder(dbType string, n int) string {
	if dbType == "postgres" {
		return fmt.Sprintf("$%d", n)
	}
	return "?"
}
