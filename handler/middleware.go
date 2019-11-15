package handler

import "net/http"

// HeadersMiddleware sets pre defined headers
func (h Handler) headersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set content type
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")

		next.ServeHTTP(w, r)
	})
}
