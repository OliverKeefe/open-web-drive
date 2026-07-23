package middleware

import (
	"net/http"
)

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Content-Disposition is not in the browser's default safe list,
		// so it must be explicitly exposed for the frontend to read it.
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")

		if r.Method == http.MethodOptions {
			w.WriteHeader((http.StatusOK))
			return
		}
		next.ServeHTTP(w, r)
	})
}
