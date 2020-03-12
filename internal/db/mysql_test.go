package db_test

import (
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/db"
)

func TestNewMySQL(t *testing.T) {
	c := config.MySQL{
		Host:     "127.0.0.1",
		Port:     3316,
		Username: "user",
		Password: "pass",
		Database: "recipes",
	}

	mdb, err := db.NewMySQL(c)
	if err != nil {
		t.Fatal(err)
	}

	result := struct{ one int }{}
	if err := mdb.QueryRow(`select 1 as one`).Scan(&result.one); err != nil {
		t.Fatal(err)
	}

	if result.one != 1 {
		t.Fatal("Invalid result")
	}
}
