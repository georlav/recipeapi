package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georlav/recipeapi/internal/config"

	"github.com/georlav/recipeapi/internal/handler"
)

func TestHandler_ContentTypeMiddleware(t *testing.T) {
	h := handler.NewHandler(nil, &config.Config{}, nil)

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
	h := handler.NewHandler(nil, &config.Config{}, nil)

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
