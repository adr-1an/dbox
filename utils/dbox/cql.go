package dbox

import (
	"os"
	"strconv"

	"github.com/gocql/gocql"
)

func openCQLSession() (*gocql.Session, error) {
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	port, _ := strconv.Atoi(portStr)

	keyspace := os.Getenv("DB_NAME")

	cluster := gocql.NewCluster(host)
	if port != 0 {
		cluster.Port = port
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: user,
			Password: os.Getenv("DB_PASS"),
		}
	}
	if keyspace != "" {
		cluster.Keyspace = keyspace
	}

	return cluster.CreateSession()
}
