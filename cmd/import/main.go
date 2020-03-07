package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/recipe"
	"github.com/georlav/recipeapi/internal/recipe/mongodb"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration, %s", err))
	}

	// initialize logger
	logger := log.New(
		os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile,
	)

	// Mongo client
	client, err := mongodb.New(cfg.Mongo)
	if err != nil {
		log.Fatalf(`mongo client error, %s`, err)
	}

	// Initialize mongo
	db := client.Database(cfg.Mongo.Database)
	rCollection := db.Collection(cfg.Mongo.RecipeCollection)
	rr := recipe.NewMongoRepo(rCollection, cfg.Mongo)

	// Import recipes from file
	file, err := os.Open("recipes.json")
	if err != nil {
		logger.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var r []struct {
			Title       string `json:"title"`
			URL         string `json:"href"`
			Ingredients string `json:"ingredients"`
			Thumbnail   string `json:"thumbnail"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &r); err != nil {
			logger.Fatal(err)
		}

		recipes := recipe.Recipes{}
		for i := range r {
			ing := strings.Split(r[i].Ingredients, ",")
			for i := range ing {
				ing[i] = strings.Trim(ing[i], " ")
			}

			recipes = append(recipes, recipe.Recipe{
				Title:       r[i].Title,
				URL:         r[i].URL,
				Ingredients: ing,
				Thumbnail:   r[i].Thumbnail,
			})
		}

		if err := rr.Insert(recipes...); err != nil {
			logger.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Fatal(err)
	}
}
