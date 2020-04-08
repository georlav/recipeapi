package handler_test

import (
	"testing"

	"github.com/georlav/recipeapi/internal/database"
	"github.com/georlav/recipeapi/internal/handler"
)

func TestEncodeEntity(t *testing.T) {
	testCase := []struct {
		desc     string
		entity   database.Recipe
		Response handler.RecipeResponseItem
	}{
		{
			desc: "Should map entity to response",
			entity: database.Recipe{
				ID:        1,
				Title:     "title1",
				URL:       "url",
				Thumbnail: "thumb.jpg",
				Ingredients: []database.Ingredient{
					{ID: 1, RecipeID: 1, Name: "sugar", CreatedAt: "2020-12-12 00:00:00", UpdatedAt: "2020-12-12 00:00:00"},
					{ID: 2, RecipeID: 1, Name: "flour", CreatedAt: "2020-12-12 00:00:00", UpdatedAt: "2020-12-12 00:00:00"},
				},
				CreatedAt: "2020-12-12 00:00:00",
				UpdatedAt: "2020-12-12 00:00:00",
			},
			Response: handler.RecipeResponseItem{},
		},
		{
			desc: "Should map entity to response",
			entity: database.Recipe{
				ID:        1,
				Thumbnail: "thumb.jpg",
				CreatedAt: "2020-12-12 00:00:00",
				UpdatedAt: "2020-12-12 00:00:00",
			},
			Response: handler.RecipeResponseItem{},
		},
	}

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.desc, func(t *testing.T) {
			if err := handler.EncodeEntity(tc.entity, &tc.Response); err != nil {
				t.Fatal(err)
			}
			if tc.entity.ID != tc.Response.ID {
				t.Fatalf("Expected name to be %d got %d", tc.entity.ID, tc.Response.ID)
			}
			if tc.entity.Title != tc.Response.Title {
				t.Fatalf("Expected title to be %s got %s", tc.entity.Title, tc.Response.Title)
			}
			if len(tc.entity.Ingredients) != len(tc.Response.Ingredients) {
				t.Fatalf("Expected ingredients length to be %d got %d", len(tc.entity.Ingredients), len(tc.Response.Ingredients))
			}
		})
	}

	t.Run("Should fail to map entity to response", func(t *testing.T) {
		if err := handler.EncodeEntity(struct{}{}, struct{}{}); err == nil {
			t.Fatal("Expected to fail due to non-pointer target")
		}
	})
}

func TestEncodeEntities(t *testing.T) {
	testCase := []struct {
		desc     string
		entity   database.Recipes
		Response handler.RecipesResponse
	}{
		{
			entity: database.Recipes{
				{
					ID:        1,
					Title:     "title1",
					URL:       "url1",
					Thumbnail: "thumb1.jpg",
					Ingredients: []database.Ingredient{
						{ID: 1, RecipeID: 1, Name: "sugar", CreatedAt: "2020-12-12 00:00:00", UpdatedAt: "2020-12-12 00:00:00"},
						{ID: 2, RecipeID: 1, Name: "flour", CreatedAt: "2020-12-12 00:00:00", UpdatedAt: "2020-12-12 00:00:00"},
					},
					CreatedAt: "2020-12-12 00:00:00",
					UpdatedAt: "2020-12-12 00:00:00",
				},
				{
					ID:        2,
					Title:     "title2",
					URL:       "url2",
					Thumbnail: "thumb2.jpg",
					CreatedAt: "2020-12-12 00:00:00",
					UpdatedAt: "2020-12-12 00:00:00",
				},
			},
			Response: handler.RecipesResponse{},
		},
		{
			entity: database.Recipes{
				{
					ID:        1,
					Title:     "title1",
					URL:       "url1",
					Thumbnail: "thumb1.jpg",
					Ingredients: []database.Ingredient{
						{ID: 1, RecipeID: 1, Name: "sugar", CreatedAt: "2020-12-12 00:00:00", UpdatedAt: "2020-12-12 00:00:00"},
						{ID: 2, RecipeID: 1, Name: "flour", CreatedAt: "2020-12-12 00:00:00", UpdatedAt: "2020-12-12 00:00:00"},
					},
					CreatedAt: "2020-12-12 00:00:00",
					UpdatedAt: "2020-12-12 00:00:00",
				},
			},
			Response: handler.RecipesResponse{},
		},
	}

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.desc, func(t *testing.T) {
			if err := handler.EncodeEntities(tc.entity, &tc.Response, "Data"); err != nil {
				t.Fatal(err)
			}

			for i, v := range *tc.Response.Data {
				if tc.entity[i].ID != v.ID {
					t.Fatalf("Expected name to be %d got %d", tc.entity[i].ID, v.ID)
				}
				if tc.entity[i].Title != v.Title {
					t.Fatalf("Expected title to be %s got %s", tc.entity[i].Title, v.Title)
				}
				if len(tc.entity[i].Ingredients) != len(v.Ingredients) {
					t.Fatalf("Expected ingredients length to be %d got %d", len(tc.entity[i].Ingredients), len(v.Ingredients))
				}
			}
		})
	}

	t.Run("Should fail to map entities due to invalid entity type", func(t *testing.T) {
		if err := handler.EncodeEntities(struct{}{}, struct{}{}, "Result"); err == nil {
			t.Fatal(err)
		}
	})

	t.Run("Should fail to map entities because response is not a pointer", func(t *testing.T) {
		if err := handler.EncodeEntities([]struct{}{}, struct{}{}, "Result"); err == nil {
			t.Fatal(err)
		}
	})

	t.Run("Should fail to map entities because response target filed is missing", func(t *testing.T) {
		if err := handler.EncodeEntities(&[]struct{}{}, &struct{}{}, "Result"); err == nil {
			t.Fatal(err)
		}
	})
}
