package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/database"
	"github.com/georlav/recipeapi/internal/handler"
	"github.com/georlav/recipeapi/internal/logger"
)

func TestHandler_User(t *testing.T) {
	testData := []struct {
		input        handler.Token
		output       handler.UserProfileResponse
		expectedCode int
	}{
		{
			handler.Token{
				UserID:   1,
				Username: "user1",
			},
			handler.UserProfileResponse{
				ID:       1,
				Username: "user1",
				FullName: "test name",
				Email:    "test@test.gr",
				Active:   true,
			},
			http.StatusOK,
		},
		{
			handler.Token{},
			handler.UserProfileResponse{},
			http.StatusUnauthorized,
		},
	}

	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		t.Fatal(err)
	}

	h := handler.NewHandler(db, cfg, logger.NewLogger(cfg.Logger))

	for i := range testData {
		tc := testData[i]

		t.Run(fmt.Sprintf(`Get profile of user with id %d`, tc.output.ID), func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("GET", "/user", nil)

			rr := httptest.NewRecorder()
			h := http.HandlerFunc(h.User)

			ctx := context.Background()
			if tc.input.UserID != 0 {
				ctx = context.WithValue(req.Context(), handler.CtxKeyToken, tc.input)
			}
			h.ServeHTTP(rr, req.WithContext(ctx))

			if rr.Code != tc.expectedCode {
				t.Fatalf("Wrong status code got %d expected %d, %s", rr.Code, tc.expectedCode, rr.Body.String())
			}
		})
	}
}

func TestHandler_SignIn(t *testing.T) {
	testData := []struct {
		input        string
		expectedCode int
	}{
		{
			`{"username": "username1", "password": "password"}`,
			http.StatusOK,
		},
		{
			`{"username": "username1", "password": "pass"}`,
			http.StatusUnauthorized,
		},
		{
			`{"username": "username", "password": "password"}`,
			http.StatusUnauthorized,
		},
		{
			`{"username": "", "password": ""}`,
			http.StatusBadRequest,
		},
		{
			`invalid input`,
			http.StatusBadRequest,
		},
	}

	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		t.Fatal(err)
	}

	h := handler.NewHandler(db, cfg, logger.NewLogger(cfg.Logger))

	for i := range testData {
		tc := testData[i]

		t.Run(`Sign in`, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/user/signin", strings.NewReader(tc.input))
			rr := httptest.NewRecorder()
			h := http.HandlerFunc(h.SignIn)
			h.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Fatalf("Wrong status code got %d expected %d, %s", rr.Code, tc.expectedCode, rr.Body.String())
			}

			if rr.Code == http.StatusOK {
				tr := handler.TokenResponse{}
				if err := json.Unmarshal(rr.Body.Bytes(), &tr); err != nil {
					t.Fatal(err)
				}

				if tr.Token == "" {
					t.Fatal("Token is empty")
				}
			}
		})
	}
}

func TestHandler_SignUp(t *testing.T) {
	testData := []struct {
		desc         string
		input        string
		expectedCode int
	}{
		{
			"Should create an account",
			`{"username":"username2","password":"password","repeatPassword":"password","fullName":"test user","email":"email@email.com"}`,
			http.StatusOK,
		},
		{
			"Should fail because passwords doesnt match",
			`{"username": "username3","password": "password3","repeatPassword": "password4",
"fullName": "full name","email": "test+2@test.com"}`,
			http.StatusBadRequest,
		},
		{
			"Should fail due to missing username",
			`{"password": "password3","repeatPassword": "password4",
"fullName": "full name","email": "test+2@test.com"}`,
			http.StatusBadRequest,
		},
		{
			"Should fail because request payload is not a valid json",
			`invalid request`,
			http.StatusBadRequest,
		},
	}

	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		t.Fatal(err)
	}

	h := handler.NewHandler(db, cfg, logger.NewLogger(cfg.Logger))

	for i := range testData {
		tc := testData[i]

		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/user/signup", strings.NewReader(tc.input))
			rr := httptest.NewRecorder()
			h := http.HandlerFunc(h.SignUp)
			h.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Fatalf("Wrong status code got %d expected %d, %s", rr.Code, tc.expectedCode, rr.Body.String())
			}

			if rr.Code == http.StatusOK {
				tr := handler.TokenResponse{}
				if err := json.Unmarshal(rr.Body.Bytes(), &tr); err != nil {
					t.Fatal(err)
				}

				if tr.Token == "" {
					t.Fatal("Token is empty")
				}
			}
		})
	}
}
