package db

import (
	"database/sql"
	"fmt"

	"github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	log = logrus.WithField("context", "db")
)

// Connect establishes a connection to the PostgreSQL database used for all SQL
// queries. This function should be called before using any other types or
// functions in the package.
func Connect(name, user, password, host string, port int) {
	d, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"dbname=%s user=%s password=%s host=%s port=%d",
			name,
			user,
			password,
			host,
			port,
		),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	db = d
}
