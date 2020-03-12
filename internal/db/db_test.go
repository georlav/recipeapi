package db_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/db"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMain(m *testing.M) {
	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Init mysql queries
	cfg.APP.Database = "MySQL"
	sqlQueries, err := db.New(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Init Mongo queries
	cfg.APP.Database = "MongoDB"
	mongoQueries, err := db.New(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Import data
	jd, err := ioutil.ReadFile("testdata/recipes.json")
	if err != nil {
		log.Fatalf("failed to load test data, %s", err)
	}
	if err := importData(sqlQueries, jd); err != nil {
		log.Fatal(err)
	}
	if err := importData(mongoQueries, jd); err != nil {
		log.Fatal(err)
	}

	status := m.Run()

	if err := tearDownMySQL(cfg.MySQL); err != nil {
		fmt.Println("MySQL clean up failed, you have to do it manually")
	}
	if err := tearDownMongo(cfg.Mongo); err != nil {
		fmt.Println("MySQL clean up failed, you have to do it manually")
	}

	os.Exit(status)
}

func importData(q db.Queryable, jsonData []byte) error {
	var data struct{ Recipes db.Recipes }
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("failed to marshal testdata, %w", err)
	}

	// Import data
	for i := range data.Recipes {
		if err := q.Insert(data.Recipes[i]); err != nil {
			return fmt.Errorf("failed to insert test data, %w", err)
		}
	}

	return nil
}

func tearDownMySQL(cfg config.MySQL) error {
	sqlDB, err := db.NewMySQL(cfg)
	if err != nil {
		return err
	}

	if _, err := sqlDB.Exec(`SET FOREIGN_KEY_CHECKS = 0`); err != nil {
		return err
	}
	if _, err := sqlDB.Exec(`TRUNCATE TABLE recipe`); err != nil {
		return err
	}
	if _, err := sqlDB.Exec(`TRUNCATE TABLE ingredient`); err != nil {
		return err
	}
	if _, err := sqlDB.Exec(`SET FOREIGN_KEY_CHECKS = 1`); err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}

	return nil
}

func tearDownMongo(cfg config.Mongo) error {
	client, err := db.NewMongoDB(cfg)
	if err != nil {
		return fmt.Errorf(`mongo client error, %s`, err)
	}
	mdb := client.Database(cfg.Database)
	col := mdb.Collection(cfg.RecipeCollection)

	if _, err = col.DeleteMany(context.Background(), bson.D{}); err != nil {
		return fmt.Errorf(`mongo mass delete error, %w`, err)
	}

	return nil
}
