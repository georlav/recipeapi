package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// HeadersMiddleware sets pre defined headers
func (h Handler) HeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set content type
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")

		next.ServeHTTP(w, r)
	})
}

// Authorization middleware assign to all routes that require users to be signed in
func (h Handler) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(auth, "Bearer ")

		tr := Token{}
		token, err := jwt.ParseWithClaims(tokenString, &tr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return []byte(h.cfg.Token.Secret), nil
		})
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "invalid token, %s"}`, err), http.StatusBadRequest)
			return
		}
		if !token.Valid {
			http.Error(w, fmt.Sprintf(`{"error": "invalid token"}`), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyToken, tr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
