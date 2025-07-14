package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pack-calculator/internal/api/handlers"
	apihttp "pack-calculator/internal/api/http"
	"pack-calculator/internal/api/middleware"
	"pack-calculator/internal/config"
	"pack-calculator/internal/domain/service"
	"pack-calculator/internal/infrastructure/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Initialize(cfg.Logging.Level, cfg.Logging.Format)
	logger.Info("Starting Pack Calculator API", map[string]interface{}{
		"version":     cfg.App.Version,
		"environment": cfg.App.Environment,
		"port":        cfg.Server.Port,
	})

	// Initialize services
	packService := service.NewPackService()
	logger.Info("Services initialized")

	// Initialize handlers
	calculationHandler := handlers.NewCalculationHandler(packService)
	healthHandler := handlers.NewHealthHandler()
	staticHandler := handlers.NewStaticHandler()
	logger.Info("Handlers initialized")

	// Initialize HTTP server with middleware
	router := apihttp.NewRouter()

	// Register routes with handler functions
	router.RegisterCalculationRoutes(calculationHandler.Calculate)
	router.RegisterHealthRoutes(healthHandler.Health, healthHandler.Ready)
	router.RegisterStaticRoutes(staticHandler.ServeUI, staticHandler.ServeStatic)
	// Wrap with logging middleware
	handler := middleware.Logging(router.Handler())

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Info("HTTP server starting", map[string]interface{}{
			"address": server.Addr,
		})

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", map[string]interface{}{
				"error": err.Error(),
			})
			os.Exit(1)
		}
	}()

	logger.Info("Server started successfully")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutdown signal received, starting graceful shutdown")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		logger.Info("Server shutdown completed")
	}
}
