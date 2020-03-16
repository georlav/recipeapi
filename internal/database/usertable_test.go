package database_test

import (
	"errors"
	"log"
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/database"
)

func TestUserTable_Get(t *testing.T) {
	testCases := []struct {
		desc   string
		input  uint64
		output *database.User
		error  error
	}{
		{
			"Should get a user",
			1,
			&database.User{
				Username: "user1",
				FullName: "test user",
				Email:    "test@test.gr",
			},
			nil,
		},
		{
			"Should fail to get a user",
			0,
			nil,
			database.ErrNoRows,
		},
	}

	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.desc, func(t *testing.T) {
			user, err := db.User.Get(tc.input)
			if err != nil && !errors.Is(err, tc.error) {
				t.Fatal(err)
			}
			if err == nil && uint64(user.ID) != tc.input {
				t.Fatalf("Invalid id, expected %d got %d", tc.input, user.ID)
			}
			if err == nil && user.Username != tc.output.Username {
				t.Fatalf("Invalid title, expected %s got %s", tc.output.Username, user.Username)
			}
		})
	}
}

func TestUserTable_GetByUsername(t *testing.T) {
	testCases := []struct {
		desc   string
		input  string
		output *database.User
		error  error
	}{
		{
			"Should get a user",
			"user1",
			&database.User{
				ID:       1,
				Username: "user1",
				FullName: "test user",
				Email:    "test@test.gr",
			},
			nil,
		},
		{
			"Should fail to get a user",
			"user999",
			nil,
			database.ErrNoRows,
		},
	}

	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.desc, func(t *testing.T) {
			user, err := db.User.GetByUsername(tc.input)
			if err != nil && !errors.Is(err, tc.error) {
				t.Fatal(err)
			}
			if err == nil && user.Username != tc.input {
				t.Fatalf("Invalid username, expected to have %s got %s", tc.input, user.Username)
			}
			if err == nil && user.Username != tc.output.Username {
				t.Fatalf("Invalid title, expected %s got %s", tc.output.Username, user.Username)
			}
		})
	}
}

func TestUserTable_Insert(t *testing.T) {
	testCases := []struct {
		desc  string
		input database.User
		error error
	}{
		{
			"Should create a new user",
			database.User{
				Username: "user2",
				FullName: "test user",
				Email:    "test2@test.gr",
			},
			nil,
		},
		{
			"Should fail to create a new user",
			database.User{
				Username: "user2",
				FullName: "test user",
				Email:    "test2@test.gr",
			},
			database.ErrDuplicateEntry,
		},
	}

	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.desc, func(t *testing.T) {
			id, err := db.User.Insert(tc.input)
			if err != nil && !errors.Is(err, tc.error) {
				t.Fatal(err)
			}
			if err == nil && id == 0 {
				t.Fatal("User expected to have an id")
			}
		})
	}
}
