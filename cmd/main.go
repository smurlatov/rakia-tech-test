package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"rakia-tech-test/internal/application/services"
	"rakia-tech-test/internal/infrastructure/loader"
	memory_repositories "rakia-tech-test/internal/infrastructure/repositories"
	"rakia-tech-test/internal/interfaces/rest"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		if logLevel, err := logrus.ParseLevel(level); err == nil {
			logger.SetLevel(logLevel)
		}
	}

	postRepo := memory_repositories.NewMemoryPostRepository()

	dataLoader := loader.NewDataLoader(postRepo, logger)
	if err := dataLoader.LoadFromFile("blog_data.json"); err != nil {
		logger.WithError(err).Warn("Failed to load initial data, starting with empty repository")
	}

	postService := services.NewPostService(postRepo, logger)

	postHandler := rest.NewPostHandler(postService, logger)

	r := rest.SetupRouter(postHandler, logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigint
		logger.WithField("signal", sig.String()).Info("Received shutdown signal")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		logger.Info("Shutting down server gracefully...")

		if err := srv.Shutdown(ctx); err != nil {
			logger.WithError(err).Error("HTTP Server Shutdown Error")
		} else {
			logger.Info("Server shutdown completed successfully")
		}
	}()

	logger.WithField("port", port).Info("Starting server")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.WithError(err).Fatal("HTTP server ListenAndServe Error")
	}

	wg.Wait()
	logger.Info("Server stopped")
}
