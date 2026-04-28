package middleware

import (
	"backend/src/internal/auth"
	"net/http"
	"strings"
)

func AuthMiddleware(a auth.Authenticator) func(http.Handler) http.Handler {
	//return func(next http.Handler) http.Handler {
	//
	//}
	panic("not implemented")
}

func Protect(a *auth.Authenticator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(h, "Bearer ")
		newCtx, ok, err := a.ValidateJWT(r.Context(), token)
		if err != nil || !ok {
			http.Error(w, "unauthorized, invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(newCtx))
	})

}
