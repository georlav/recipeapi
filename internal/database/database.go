package database

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"github.com/georlav/recipeapi/internal/config"
)

type Database struct {
	Handle     *sql.DB
	Recipe     *RecipeTable
	Ingredient *IngredientTable
	User       *UserTable
}

func New(c config.Database) (*Database, error) {
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

	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetMaxOpenConns(c.MaxOpenConns)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping error, %w", err)
	}

	return &Database{
		Handle:     db,
		Recipe:     NewRecipeTable(db),
		Ingredient: NewIngredientTable(db),
		User:       NewUserTable(db),
	}, nil
}
