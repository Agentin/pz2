package middleware

import (
	"net/http"
	"strings"

	"github.com/student/tech-ip-sem2/services/tasks/client/authclient"
	"github.com/student/tech-ip-sem2/shared/middleware"
)

// AuthMiddleware проверяет токен через Auth service.
func AuthMiddleware(authClient *authclient.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}
			token := parts[1]

			requestID := middleware.GetRequestID(r.Context())
			valid, err := authClient.Verify(r.Context(), token, requestID)
			if err != nil {
				// Auth недоступен или другая ошибка – возвращаем 500 (fail closed)
				http.Error(w, `{"error":"authorization service unavailable"}`, http.StatusInternalServerError)
				return
			}
			if !valid {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			// токен валиден – продолжаем
			next.ServeHTTP(w, r)
		})
	}
}
