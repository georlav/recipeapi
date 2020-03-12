package db

import (
	"database/sql"
	"fmt"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/go-sql-driver/mysql"
)

// New returns a new db handle
func NewMySQL(c config.MySQL) (*sql.DB, error) {
	dsn, err := mysql.ParseDSN(
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.Username, c.Password, c.Host, c.Port, c.Database),
	)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("open error, %w", err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping error, %w", err)
	}

	return db, nil
}
