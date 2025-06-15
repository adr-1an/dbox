package dbox

import (
	"database/sql"
	"fmt"

	"github.com/gocql/gocql"
)

var SharedDB *sql.DB
var SharedSession *gocql.Session

func sqlConn(dbType, dsn string) (*sql.DB, func(), error) {
	if SharedDB != nil {
		return SharedDB, func() {}, nil
	}
	db, err := sql.Open(dbType, dsn)
	return db, func() { db.Close() }, err
}

func cqlConn() (*gocql.Session, func(), error) {
	if SharedSession != nil {
		return SharedSession, func() {}, nil
	}
	sess, err := OpenCQLSession()
	return sess, func() { sess.Close() }, err
}

func placeholder(dbType string, n int) string {
	if dbType == "postgres" {
		return fmt.Sprintf("$%d", n)
	}
	return "?"
}
