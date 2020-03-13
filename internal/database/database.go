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
}

func New(c config.MySQL) (*Database, error) {
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

	return &Database{
		Handle:     db,
		Recipe:     NewRecipeTable(db),
		Ingredient: NewIngredientTable(db),
	}, nil
}
