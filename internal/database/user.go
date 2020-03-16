package database

// User entity
type User struct {
	ID        int64
	Username  string
	Password  string `json:"-"`
	FullName  string
	Email     string
	Active    bool
	CreatedAt string
	UpdatedAt string
}

// Users slice or user entities
type Users []User
