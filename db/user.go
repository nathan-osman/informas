package db

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User represents a registered user on the website.
type User struct {
	ID         int
	Username   string
	Password   string
	Email      string
	IsAdmin    bool
	IsDisabled bool
}

// migrateUsersTable executes the SQL necessary to create the Users table.
func migrateUsersTable(t *Token) error {
	_, err := t.exec(
		`
        CREATE TABLE IF NOT EXISTS Users (
            ID         SERIAL PRIMARY KEY,
            Username   VARCHAR(40) NOT NULL UNIQUE,
            Password   VARCHAR(80) NOT NULL,
            Email      VARCHAR(100),
            IsAdmin    BOOLEAN,
            IsDisabled BOOLEAN
        )
        `,
	)
	return err
}

// AllUsers retrieves all registered users.
func AllUsers(t *Token, sort string) ([]*User, error) {
	r, err := t.query(
		`
        SELECT ID, Username, Password, Email, IsAdmin, IsDisabled
        FROM Users ORDER BY $1
        `,
		sort,
	)
	if err != nil {
		return nil, err
	}
	users := make([]*User, 0, 1)
	for r.Next() {
		u := &User{}
		if err := r.Scan(
			&u.ID,
			&u.Username,
			&u.Password,
			&u.Email,
			&u.IsAdmin,
			&u.IsDisabled,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// FindUser attempts to retrieve a user using the specified field.
func FindUser(t *Token, field string, value interface{}) (*User, error) {
	u := &User{}
	err := t.queryRow(
		fmt.Sprintf(
			`
            SELECT ID, Username, Password, Email, IsAdmin, IsDisabled
            FROM Users WHERE %s = $1
            `,
			field,
		),
		value,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.Email,
		&u.IsAdmin,
		&u.IsDisabled,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Save updates the object in the database. If the user ID is set to 0, a new
// user is instead created and their ID set to 0.
func (u *User) Save(t *Token) error {
	if u.ID == 0 {
		var id int
		err := t.queryRow(
			`
            INSERT INTO Users (Username, Password, Email, IsAdmin, IsDisabled)
            VALUES ($1, $2, $3, $4, $5) RETURNING ID
            `,
			u.Username,
			u.Password,
			u.Email,
			u.IsAdmin,
			u.IsDisabled,
		).Scan(&id)
		if err != nil {
			return err
		}
		u.ID = id
		return nil
	} else {
		_, err := t.exec(
			`
            UPDATE Users SET Username=$1, Password=$2, Email=$3, IsAdmin=$4, IsDisabled=$5
            WHERE ID = $6
            `,
			u.Username,
			u.Password,
			u.Email,
			u.IsAdmin,
			u.IsDisabled,
			u.ID,
		)
		return err
	}
}

// Delete completely destroys the user and all data associated with them.
func (u *User) Delete(t *Token) error {
	_, err := t.exec(
		`
        DELETE FROM Users WHERE ID = $1
        `,
		u.ID,
	)
	return err
}

// Authenticate hashes the provided password and compares it to the value
// stored in the database.
func (u *User) Authenticate(password string) error {
	h, err := base64.StdEncoding.DecodeString(u.Password)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(h, []byte(password))
}

// SetPassword hashes the provided password and assigns it to the user. Note
// that this does not update the database - use Save() to do that.
func (u *User) SetPassword(password string) error {
	h, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}
	u.Password = base64.StdEncoding.EncodeToString(h)
	return nil
}
