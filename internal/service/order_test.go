package service

import (
	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/entities"
	"L0/internal/repository"
	mock_repository "L0/internal/repository/mocks"
	"L0/pkg/db"
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

type mockLogger struct{}

func NewMockLogger() *mockLogger {
	return &mockLogger{}
}

func (m *mockLogger) Info(msg string, args ...any)  {}
func (m *mockLogger) Error(msg string, args ...any) {}
func (m *mockLogger) Debug(msg string, args ...any) {}

func BenchmarkGetOrderByUID_CacheHit(b *testing.B) {
	mockRepo := mock_repository.NewMockOrderRepository()
	cache := cache.NewLRUCache(3)
	logger := NewMockLogger()
	service := NewOrder(mockRepo, cache, logger)

	orderUID := "test-order-uid"
	order := &entities.Order{OrderUID: orderUID}
	cache.Set(order)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetOrderByUID(ctx, orderUID)
	}
}

func BenchmarkGetOrderByUID_DB(b *testing.B) {
	ctx := context.Background()

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	cfg := &config.Config{DB: config.DB{
		Host:     "localhost",
		Port:     "5432",
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSL:      "disable",
	}}

	db, err := db.NewPostgresDB(&cfg.DB)
	if err != nil {
		b.Fatalf("failed to initialize postgres: %s", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	cache := cache.NewLRUCache(0)
	logger := NewMockLogger()
	service := NewService(repo, cache, logger)

	orderUID := "test-order-uid"
	order := &entities.Order{OrderUID: orderUID}
	_ = service.Order.SaveOrder(ctx, order)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.Order.GetOrderByUID(ctx, orderUID)
	}
}
