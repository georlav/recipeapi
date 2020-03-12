package db_test

import (
	"log"
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/db"
)

func TestNewMySQL(t *testing.T) {
	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}

	mdb, err := db.NewMySQL(cfg.MySQL)
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
