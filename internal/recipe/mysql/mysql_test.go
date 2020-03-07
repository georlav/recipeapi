package mysql_test

import (
	"testing"

	"github.com/georlav/recipeapi/internal/config"

	"github.com/georlav/recipeapi/internal/recipe/mysql"
)

func TestNew(t *testing.T) {
	c := config.MySQL{
		Host:     "127.0.0.1",
		Port:     3316,
		Username: "user",
		Password: "pass",
		Database: "test",
	}

	db, err := mysql.New(c)
	if err != nil {
		t.Fatal(err)
	}

	result := struct{ one int }{}
	if err := db.QueryRow(`select 1 as one`).Scan(&result.one); err != nil {
		t.Fatal(err)
	}

	if result.one != 1 {
		t.Fatal("Invalid result")
	}
}
