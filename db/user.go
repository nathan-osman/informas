package db

import (
	"encoding/base64"

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
func migrateUsersTable() error {
	_, err := db.Exec(
		`
        CREATE TABLE IF NOT EXISTS Users (
            ID         SERIAL PRIMARY KEY,
            Username   VARCHAR(40) NOT NULL,
            Password   VARCHAR(80) NOT NULL,
            Email      VARCHAR(100),
            IsAdmin    BOOLEAN,
            IsDisabled BOOLEAN
        )
        `,
	)
	return err
}

// FindUser attempts to retrieve a user using the specified field.
func FindUser(field string, value interface{}) (*User, error) {
	u := &User{}
	err := db.QueryRow(
		`
        SELECT ID, Username, Password, Email, IsAdmin, IsDisabled
        FROM Users WHERE $1 = $2
        `,
		field,
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
func (u *User) Save() error {
	if u.ID == 0 {
		var id int
		err := db.QueryRow(
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
		_, err := db.Exec(
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
func (u *User) Delete() error {
	_, err := db.Exec(
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
