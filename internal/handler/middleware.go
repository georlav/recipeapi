package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// HeadersMiddleware sets pre defined headers
func (h Handler) ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set content type
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		next.ServeHTTP(w, r)
	})
}

// HeadersMiddleware sets pre defined headers
func (h Handler) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cross Origin Resource Sharing
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Max-Age", "86400")

		next.ServeHTTP(w, r)
	})
}

// Authorization middleware assign to all routes that require users to be signed in
func (h Handler) AuthorizationMiddleware(next http.Handler) http.Handler {
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
		if err != nil || !token.Valid {
			h.respondError(w, APIError{
				Message:    "invalid token",
				StatusCode: http.StatusUnauthorized,
			})
			return
		}

		ctx := context.WithValue(r.Context(), CtxKeyToken, tr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
