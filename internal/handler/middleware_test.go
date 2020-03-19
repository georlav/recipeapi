package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/georlav/recipeapi/internal/logger"

	"github.com/dgrijalva/jwt-go"
	"github.com/georlav/recipeapi/internal/config"

	"github.com/georlav/recipeapi/internal/handler"
)

func TestHandler_Authorization(t *testing.T) {
	cfg, err := config.Load("testdata/config.json")
	if err != nil {
		t.Fatal(err)
	}

	token1 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uname": "user1",
		"uid":   1,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Duration(5) * time.Minute).Unix(),
	})
	tokenSigned1, err := token1.SignedString([]byte(cfg.Token.Secret))
	if err != nil {
		t.Fatal("signature error, %w", err)
	}

	token2 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uname": "user1",
		"uid":   1,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Duration(-60) * time.Minute).Unix(),
	})
	tokenSigned2, err := token2.SignedString([]byte(cfg.Token.Secret))
	if err != nil {
		t.Fatal("signature error, %w", err)
	}

	testCases := []struct {
		desc         string
		token        string
		expectedCode int
	}{
		{
			"Token should be valid",
			tokenSigned1,
			http.StatusNoContent,
		},
		{
			"Token should not be valid",
			"xxx.yyyy.zzzz",
			http.StatusUnauthorized,
		},
		{
			"Token should not be valid because it is expired",
			tokenSigned2,
			http.StatusUnauthorized,
		},
	}

	h := handler.NewHandler(nil, *cfg, logger.NewLogger(cfg.Logger))

	for i := range testCases {
		tc := testCases[i]

		t.Run(`Sign in`, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tc.token))
			rr := httptest.NewRecorder()

			h.AuthorizationMiddleware(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNoContent)
				}),
			).ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Fatalf("Wrong status code got %d expected %d, %s", rr.Code, tc.expectedCode, rr.Body.String())
			}
		})
	}
}

func TestHandler_ContentTypeMiddleware(t *testing.T) {
	h := handler.NewHandler(nil, config.Config{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()

	h.ContentTypeMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	).ServeHTTP(rr, req)

	if rr.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatal("Invalid content type")
	}
}

func TestHandler_CorsMiddleware(t *testing.T) {
	h := handler.NewHandler(nil, config.Config{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()

	h.CorsMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	).ServeHTTP(rr, req)

	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("Invalid origin")
	}
	if rr.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
		t.Fatal("Invalid content type")
	}
}
