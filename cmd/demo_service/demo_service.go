package main

import (
	"L0/internal/config"
	delivery "L0/internal/delivery/http"
	"L0/internal/delivery/kafka"
	"L0/internal/repository"
	"L0/internal/server"
	"L0/internal/service"
	"L0/pkg/db"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("initializing server", slog.String("address", cfg.HTTPServer.Address))

	db, err := db.NewPostgresDB(&cfg.DB)
	if err != nil {
		log.Error("failed to initialize postgres", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := delivery.NewHandler(service)

	consumer := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		cfg.GroupID,
		service.Order,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Error("kafka consumer error", slog.String("error", err.Error()))
		}
	}()
	defer consumer.Close()

	srv := server.NewServer(&cfg.HTTPServer, handler)
	go func() {
		if err := srv.Start(); err != nil {
			log.Error("server error", slog.String("error", err.Error()))
			cancel()
		}
		// log.Info("server started", slog.String("address", cfg.HTTPServer.Address))
	}()
	// if err := srv.Start(); err != nil {
	// 	log.Error("server error", slog.String("error", err.Error()))
	// }
	log.Info("server started", slog.String("address", cfg.HTTPServer.Address))
	<-sigs
	log.Info("shutting down gracefully")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
