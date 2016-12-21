package db

import (
	"database/sql"
)

// Token abstracts access to the database, allowing models to use both a
// database connection and a transaction in the same way.
type Token struct {
	tx *sql.Tx
}

func (t *Token) query(query string, args ...interface{}) (*sql.Rows, error) {
	if t.tx != nil {
		return t.tx.Query(query, args...)
	}
	return db.Query(query, args...)
}

func (t *Token) queryRow(query string, args ...interface{}) *sql.Row {
	if t.tx != nil {
		return t.tx.QueryRow(query, args...)
	}
	return db.QueryRow(query, args...)
}

func (t *Token) exec(query string, args ...interface{}) (sql.Result, error) {
	if t.tx != nil {
		return t.tx.Exec(query, args...)
	}
	return db.Exec(query, args...)
}
