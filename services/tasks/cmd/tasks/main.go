package main

import (
	"context"
	"log/slog"
	"net/http" // стандартный пакет
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/student/tech-ip-sem2/services/tasks/internal/grpcclient"
	taskshttp "github.com/student/tech-ip-sem2/services/tasks/internal/http" // алиас для внутреннего пакета
	"github.com/student/tech-ip-sem2/services/tasks/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	tasksPort := os.Getenv("TASKS_PORT")
	if tasksPort == "" {
		tasksPort = "8082"
	}
	authGrpcAddr := os.Getenv("AUTH_GRPC_ADDR")
	if authGrpcAddr == "" {
		authGrpcAddr = "localhost:50051"
	}

	taskService := service.NewTaskService()

	// Создаём gRPC клиент
	authClient, err := grpcclient.NewAuthClient(authGrpcAddr)
	if err != nil {
		logger.Error("failed to create auth gRPC client", "error", err)
		os.Exit(1)
	}
	defer authClient.Close()

	// Используем алиас taskshttp
	router := taskshttp.NewRouter(taskService, authClient, logger)

	server := &http.Server{
		Addr:         ":" + tasksPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Tasks service started", "port", tasksPort, "auth_grpc_addr", authGrpcAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()

	<-done
	logger.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("shutdown error", "error", err)
	}
}
