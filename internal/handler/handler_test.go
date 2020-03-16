package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	// load config
	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}

	// load data to import
	jd, err := ioutil.ReadFile("testdata/recipes.json")
	if err != nil {
		log.Fatalf("failed to load test data, %s", err)
	}

	// Create recipes from data
	var data struct{ Recipes database.Recipes }
	if err := json.Unmarshal(jd, &data); err != nil {
		log.Fatalf("failed to marshal testdata, %s", err)
	}

	// Get a recipe handle
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

	// Create a user
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.User.Insert(database.User{
		Username: "username1",
		Password: string(hash),
		FullName: "test user",
		Email:    "test@test.gr",
	}); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if _, err := db.Handle.Exec(`SET FOREIGN_KEY_CHECKS = 0`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Handle.Exec(`TRUNCATE TABLE user`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Handle.Exec(`TRUNCATE TABLE recipe`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Handle.Exec(`TRUNCATE TABLE ingredient`); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Handle.Exec(`SET FOREIGN_KEY_CHECKS = 1`); err != nil {
		log.Fatal(err)
	}
	if err := db.Handle.Close(); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}
