package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/student/tech-ip-sem2/services/tasks/internal/grpcclient"
	"github.com/student/tech-ip-sem2/shared/middleware"
)

func AuthMiddleware(authClient *grpcclient.AuthClient) func(http.Handler) http.Handler {
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

			// Получаем request-id из контекста (установлен ранее middleware.RequestIDMiddleware)
			requestID := middleware.GetRequestID(r.Context())

			// Создаём контекст с дедлайном 2 секунды
			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()

			// Прокидываем request-id в gRPC metadata
			ctx = metadata.AppendToOutgoingContext(ctx, "x-request-id", requestID)

			valid, _, err := authClient.Verify(ctx, token)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					http.Error(w, `{"error":"auth service timeout"}`, http.StatusGatewayTimeout) // 504
					return
				}
				// Недоступность или другая ошибка
				http.Error(w, `{"error":"authorization service unavailable"}`, http.StatusServiceUnavailable) // 503
				return
			}
			if !valid {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
