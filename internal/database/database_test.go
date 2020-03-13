package database_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/database"
)

func TestMain(m *testing.M) {
	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile("testdata/recipes.json")
	if err != nil {
		log.Fatalf("failed to load test data, %s", err)
	}

	var data struct{ Recipes database.Recipes }
	if err := json.Unmarshal(b, &data); err != nil {
		log.Fatalf("failed to unmarshal test data, %s", err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	// Import data
	for i := range data.Recipes {
		if _, err := db.Recipe.Insert(data.Recipes[i]); err != nil {
			log.Fatalf("failed to insert test data, %s", err)
		}
	}

	status := m.Run()

	if _, err := db.Handle.Exec(`SET FOREIGN_KEY_CHECKS = 0`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Handle.Exec(`TRUNCATE TABLE ingredient`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Handle.Exec(`TRUNCATE TABLE recipe`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Handle.Exec(`SET FOREIGN_KEY_CHECKS = 1`); err != nil {
		log.Fatal(err)
	}
	if err := db.Handle.Close(); err != nil {
		log.Fatal(err)
	}

	os.Exit(status)
}

func TestNewDatabase(t *testing.T) {
	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		t.Fatal(err)
	}

	result := struct{ one int }{}
	if err := db.Handle.QueryRow(`select 1 as one`).Scan(&result.one); err != nil {
		t.Fatal(err)
	}

	if result.one != 1 {
		t.Fatal("Invalid result")
	}
}
