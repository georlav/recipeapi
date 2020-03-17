package database

import (
	"database/sql"
	"fmt"
	"strings"
)

const userColumns = "u.id, u.username, u.fullName, u.email, u.active, u.created_at, u.updated_at"

// UserTable object
type UserTable struct {
	db   *sql.DB
	name string
}

// NewUserTable object
func NewUserTable(db *sql.DB) *UserTable {
	return &UserTable{
		db:   db,
		name: "user u",
	}
}

// Get user by id
func (ut *UserTable) Get(id uint64) (*User, error) {
	// nolint:gosec
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE id = ?`, userColumns, ut.name)

	var u User
	if err := ut.db.QueryRow(query, id).Scan(
		&u.ID, &u.Username, &u.FullName, &u.Email, &u.Active, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &u, nil
}

// Get user by id
func (ut *UserTable) GetByUsername(uName string) (*User, error) {
	// nolint:gosec
	query := fmt.Sprintf(`SELECT %s, password FROM %s WHERE username = ?`, userColumns, ut.name)

	var u User
	if err := ut.db.QueryRow(query, uName).Scan(
		&u.ID, &u.Username, &u.FullName, &u.Email, &u.Active, &u.CreatedAt, &u.UpdatedAt, &u.Password,
	); err != nil {
		return nil, err
	}

	return &u, nil
}

// Insert a new user and return a unique identifier
func (ut *UserTable) Insert(u User) (int64, error) {
	q := `INSERT INTO user (username, password, fullName, email) VALUES (?, ?, ?, ?)`
	res, err := ut.db.Exec(q, u.Username, u.Password, u.FullName, u.Email)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			return 0, ErrDuplicateEntry
		}
		return 0, fmt.Errorf("user error, %w", err)
	}

	var uID int64
	if uID, err = res.LastInsertId(); err != nil {
		return 0, fmt.Errorf("user error, %w", err)
	}

	return uID, nil
}
