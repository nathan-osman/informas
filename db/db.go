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

// Transaction begins a new transaction and passes it to the provided callback.
// If no error is returned, Commit() is invoked - otherwise, Rollback(). The
// error is returned.
func Transaction(f func(*Token) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	t := &Token{tx: tx}
	if err := f(t); err != nil {
		t.tx.Rollback()
		return err
	}
	t.tx.Commit()
	return nil
}

// Migrate performs all database migrations.
func Migrate() error {
	tableMigrations := []func(*Token) error{
		migrateConfigTable,
		migrateUsersTable,
	}
	return Transaction(func(t *Token) error {
		for _, f := range tableMigrations {
			if err := f(t); err != nil {
				return err
			}
		}
		return nil
	})
}
