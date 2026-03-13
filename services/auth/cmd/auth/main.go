package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	authv1 "github.com/student/tech-ip-sem2/pkg/api/auth/v1"
	authgrpc "github.com/student/tech-ip-sem2/services/auth/internal/grpc"
	authhttp "github.com/student/tech-ip-sem2/services/auth/internal/http"
	"github.com/student/tech-ip-sem2/services/auth/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// HTTP порт
	httpPort := os.Getenv("AUTH_PORT")
	if httpPort == "" {
		httpPort = "8081"
	}
	// gRPC порт
	grpcPort := os.Getenv("AUTH_GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	svc := service.NewAuthService()

	// ----- gRPC сервер -----
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		logger.Error("failed to listen grpc", "error", err)
		os.Exit(1)
	}
	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, authgrpc.NewAuthServer(svc))

	go func() {
		logger.Info("gRPC server started", "port", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("gRPC server error", "error", err)
		}
	}()

	// ----- HTTP сервер -----
	router := authhttp.NewRouter(svc, logger)
	httpServer := &http.Server{
		Addr:         ":" + httpPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Info("HTTP server started", "port", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", "error", err)
		}
	}()

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	logger.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcServer.GracefulStop()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("HTTP shutdown error", "error", err)
	}
}
