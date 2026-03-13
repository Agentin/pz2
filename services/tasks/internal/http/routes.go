package http

import (
	"log/slog"
	"net/http"

	"github.com/student/tech-ip-sem2/services/tasks/internal/grpcclient"
	"github.com/student/tech-ip-sem2/services/tasks/internal/http/handlers"
	authMiddleware "github.com/student/tech-ip-sem2/services/tasks/internal/http/handlers/middleware"
	"github.com/student/tech-ip-sem2/services/tasks/internal/service"
	"github.com/student/tech-ip-sem2/shared/middleware"
)

func NewRouter(taskService *service.TaskService, authClient *grpcclient.AuthClient, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	protected := authMiddleware.AuthMiddleware(authClient)

	mux.Handle("POST /v1/tasks", protected(http.HandlerFunc(handlers.CreateTaskHandler(taskService))))
	mux.Handle("GET /v1/tasks", protected(http.HandlerFunc(handlers.GetTasksHandler(taskService))))
	mux.Handle("GET /v1/tasks/{id}", protected(http.HandlerFunc(handlers.GetTaskHandler(taskService))))
	mux.Handle("PATCH /v1/tasks/{id}", protected(http.HandlerFunc(handlers.UpdateTaskHandler(taskService))))
	mux.Handle("DELETE /v1/tasks/{id}", protected(http.HandlerFunc(handlers.DeleteTaskHandler(taskService))))

	handler := middleware.RequestIDMiddleware(mux)
	handler = middleware.LoggingMiddleware(logger)(handler)
	return handler
}
