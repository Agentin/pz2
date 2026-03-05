package http

import (
	"log/slog"
	"net/http"

	"github.com/student/tech-ip-sem2/services/auth/internal/http/handlers"
	"github.com/student/tech-ip-sem2/services/auth/internal/service"
	"github.com/student/tech-ip-sem2/shared/middleware"
)

func NewRouter(svc *service.AuthService, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/auth/login", handlers.LoginHandler(svc))
	mux.HandleFunc("GET /v1/auth/verify", handlers.VerifyHandler(svc))

	// Оборачиваем в middleware
	handler := middleware.RequestIDMiddleware(mux)
	handler = middleware.LoggingMiddleware(logger)(handler)

	return handler
}
