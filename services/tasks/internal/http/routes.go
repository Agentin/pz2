package http

import (
	"log/slog"
	"net/http"

	"github.com/student/tech-ip-sem2/services/tasks/client/authclient"
	"github.com/student/tech-ip-sem2/services/tasks/internal/http/handlers"
	authMiddleware "github.com/student/tech-ip-sem2/services/tasks/internal/http/handlers/middleware"
	"github.com/student/tech-ip-sem2/services/tasks/internal/service"
	"github.com/student/tech-ip-sem2/shared/middleware"
)

func NewRouter(taskService *service.TaskService, authClient *authclient.AuthClient, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	// Защищённые эндпоинты (все, кроме публичных, если бы они были)
	protected := authMiddleware.AuthMiddleware(authClient)

	// Регистрируем маршруты
	mux.Handle("POST /v1/tasks", protected(http.HandlerFunc(handlers.CreateTaskHandler(taskService))))
	mux.Handle("GET /v1/tasks", protected(http.HandlerFunc(handlers.GetTasksHandler(taskService))))
	mux.Handle("GET /v1/tasks/{id}", protected(http.HandlerFunc(handlers.GetTaskHandler(taskService))))
	mux.Handle("PATCH /v1/tasks/{id}", protected(http.HandlerFunc(handlers.UpdateTaskHandler(taskService))))
	mux.Handle("DELETE /v1/tasks/{id}", protected(http.HandlerFunc(handlers.DeleteTaskHandler(taskService))))

	// Глобальные middleware
	handler := middleware.RequestIDMiddleware(mux)
	handler = middleware.LoggingMiddleware(logger)(handler)

	return handler
}
