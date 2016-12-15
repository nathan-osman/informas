package models

import (
	"golang.org/x/crypto/bcrypt"

	"encoding/base64"
)

// User represents a registered user on the website.
type User struct {
	ID       int
	Username string
	Password string
	Email    string
	IsAdmin  bool
	IsActive bool
}

// CreateUserTable executes the SQL necessary to create the User table.
func CreateUserTable() error {
	_, err := db.Exec(
		`
        CREATE TABLE User (
            ID SERIAL PRIMARY KEY,
            Username VARCHAR(40) NOT NULL,
            Password VARCHAR(80) NOT NULL,
            Email    VARCHAR(100),
            IsAdmin  BOOLEAN,
            IsActive BOOLEAN
        )
        `,
	)
	return err
}

// GetUser retrieves a User object from the database that matches the given
// user ID. If the ID does not match any users, ErrUserDoesNotExist is
// returned.
func GetUser(userID int) (*User, error) {
	u := &User{}
	err := db.QueryRow(
		`
        SELECT ID, Username, Password, Email, IsAdmin, IsActive
        FROM User WHERE ID = ?
        `,
		userID,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.Email,
		&u.IsAdmin,
		&u.IsActive,
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
		r, err := db.Exec(
			`
            INSERT INTO User (Username, Password, Email, IsAdmin, IsActive)
            VALUES (?, ?, ?, ?, ?);
            `,
			u.Username,
			u.Password,
			u.Email,
			u.IsAdmin,
			u.IsActive,
		)
		if err != nil {
			return err
		}
		i, err := r.LastInsertId()
		if err != nil {
			return err
		}
		u.ID = int(i)
		return nil
	} else {
		_, err := db.Exec(
			`
            UPDATE User SET Username=?, Password=?, Email=?, IsAdmin=?, IsActive=?
            WHERE ID = ?
            `,
			u.Username,
			u.Password,
			u.Email,
			u.IsAdmin,
			u.IsActive,
			u.ID,
		)
		return err
	}
}

// Delete completely destroys the user and all data associated with them.
func (u *User) Delete() error {
	_, err := db.Exec(
		`
        DELETE FROM User WHERE ID = ?
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
