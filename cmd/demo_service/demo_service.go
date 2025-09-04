package main

import (
	"L0/internal/cache"
	"L0/internal/config"
	delivery "L0/internal/delivery/http"
	"L0/internal/delivery/kafka"
	"L0/internal/repository"
	"L0/internal/server"
	"L0/internal/service"
	"L0/pkg/db"
	"L0/pkg/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	logger := logger.NewSlogLogger(cfg.Env)

	db, err := db.NewPostgresDB(&cfg.DB)
	if err != nil {
		logger.Error("failed to initialize postgres", "error", err.Error())
		return
	}
	defer db.Close()
	// logger.Info("postgres connected", slog.String("host", cfg.DB.Host), slog.String("port", cfg.DB.Port))
	logger.Info("postgres connected", "host", cfg.DB.Host, "port", cfg.DB.Port)

	repository := repository.NewRepository(db)
	cache := cache.NewLRUCache(cfg.Cache.Capacity)
	service := service.NewService(repository, cache, logger)

	ctx := context.Background()
	if err := service.Order.LoadCache(ctx); err != nil {
		// logger.Error("failed to load cache", slog.String("error", err.Error()))
		logger.Error("failed to load cache", "error", err.Error())
		return
	}
	// logger.Info("cache loader successfully")

	handler := delivery.NewHandler(service, logger)

	consumer := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		cfg.GroupID,
		service.Order,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := consumer.Start(ctx); err != nil {
			// logger.Error("kafka consumer error", slog.String("error", err.Error()))
			logger.Error("kafka consumer error", "error", err.Error())
		}
	}()
	defer consumer.Close()

	srv := server.NewServer(&cfg.HTTPServer, handler)
	go func() {
		if err := srv.Start(); err != nil {
			// logger.Error("server error", slog.String("error", err.Error()))
			logger.Error("server error", "error", err.Error())
			cancel()
		}
	}()
	defer srv.Shutdown(ctx)

	// logger.Info("server started", slog.String("address", cfg.HTTPServer.Address))
	logger.Info("server started", "address", cfg.HTTPServer.Address)
	<-sigs
	logger.Info("shutting down gracefully")
}
