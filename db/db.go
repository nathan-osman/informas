package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Connect establishes a connection to the PostgreSQL database used for all SQL
// queries. This function should be called before using any other types or
// functions in the package.
func Connect(name, user, password, host string, port int) error {
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
		return err
	}
	db = d
	return nil
}

// Migrate performs all database migrations.
func Migrate() error {
	tableMigrations := []func() error{
		migrateUsersTable,
	}
	for _, f := range tableMigrations {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
