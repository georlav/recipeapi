package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

	h := handler.NewHandler(db, *cfg, &logger.Logger{})

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
		input        handler.SignInRequest
		expectedCode int
	}{
		{
			handler.SignInRequest{
				Username: "username1",
				Password: "password",
			},
			http.StatusOK,
		},
		{
			handler.SignInRequest{
				Username: "username1",
				Password: "pass",
			},
			http.StatusUnauthorized,
		},
		{
			handler.SignInRequest{
				Username: "username",
				Password: "password",
			},
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

	h := handler.NewHandler(db, config.Config{}, &logger.Logger{})

	for i := range testData {
		tc := testData[i]

		t.Run(fmt.Sprintf(`Sign in as %s`, tc.input.Username), func(t *testing.T) {
			t.Parallel()

			b, err := json.Marshal(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/user/signin", bytes.NewReader(b))

			// initialize response recorder to monitor handler response data
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
		input        handler.SignUpRequest
		expectedCode int
	}{
		{
			handler.SignUpRequest{
				Email:          "test+1@test.com",
				FullName:       "full name",
				Username:       "username2",
				Password:       "password2",
				RepeatPassword: "password2",
			},
			http.StatusOK,
		},
		{
			handler.SignUpRequest{
				Email:          "test+2@test.com",
				FullName:       "full name",
				Username:       "username3",
				Password:       "password3",
				RepeatPassword: "password4",
			},
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

	h := handler.NewHandler(db, config.Config{}, &logger.Logger{})

	for i := range testData {
		tc := testData[i]

		t.Run(fmt.Sprintf(`Sign in as %s`, tc.input.Username), func(t *testing.T) {
			t.Parallel()

			b, err := json.Marshal(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/user/signup", bytes.NewReader(b))

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
