package main

import (
	"context"
	"fmt"
	"gopher-identity-service/internal/presentation/http/handlers/auth"
	"gopher-identity-service/pkg/jwt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"gopher-identity-service/internal/application/usecases"
	"gopher-identity-service/internal/config"
	"gopher-identity-service/internal/infrastructure/database"
	"gopher-identity-service/internal/infrastructure/database/repositories"
	grpcServer "gopher-identity-service/internal/presentation/grpc"
	httpRouter "gopher-identity-service/internal/presentation/http"
	"gopher-identity-service/internal/presentation/http/handlers/user"
	"gopher-identity-service/pkg/logger"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// Core dependencies
	container.Provide(config.LoadConfig)
	container.Provide(logger.NewLogger)
	container.Provide(func(cfg *config.Config) *config.JWTConfig {
		return &cfg.JWT
	})

	// Infrastructure
	container.Provide(database.NewPostgresDB)
	container.Provide(repositories.NewUserPostgresRepository)

	// Pkg
	container.Provide(jwt.NewManager)

	// Application
	container.Provide(usecases.NewSSOUseCase)
	container.Provide(usecases.NewSignUpUseCase)

	// Presentation
	container.Provide(user.NewGetProfileHandler)
	container.Provide(user.NewRouter)

	container.Provide(auth.NewSignUpHandler)
	container.Provide(auth.NewRouter)

	container.Provide(httpRouter.NewRouter)
	container.Provide(grpcServer.NewGRPCServer)

	return container
}

func main() {
	container := BuildContainer()

	err := container.Invoke(func(cfg *config.Config, log *zap.Logger, router *gin.Engine, grpcSrv *grpc.Server) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		// Start HTTP Server
		httpServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.App.HTTPPort),
			Handler: router,
		}

		go func() {
			log.Info("Starting HTTP Server", zap.Int("port", cfg.App.HTTPPort))
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal("HTTP Server failed to start", zap.Error(err))
			}
		}()

		// Start gRPC Server
		go func() {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.App.GRPCPort))
			if err != nil {
				log.Fatal("Failed to listen for gRPC", zap.Error(err))
			}
			log.Info("Starting gRPC Server", zap.Int("port", cfg.App.GRPCPort))
			if err := grpcSrv.Serve(lis); err != nil {
				log.Fatal("gRPC Server failed to start", zap.Error(err))
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the servers
		<-ctx.Done()
		log.Info("Shutting down gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error("HTTP Server forced to shutdown", zap.Error(err))
		}

		grpcSrv.GracefulStop()

		log.Info("Servers exited gracefully")
	})

	if err != nil {
		fmt.Printf("Error starting application: %v\n", err)
		os.Exit(1)
	}
}
