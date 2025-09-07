package main

import (
	"L0/internal/entities"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "orders"

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	log.Println("Producer started")

	for {
		order := generateOrder()
		data, err := json.Marshal(order)
		if err != nil {
			log.Printf("failed to marshal order: %v", err)
			continue
		}

		msg := kafka.Message{
			Key:   []byte(order.OrderUID),
			Value: data,
		}

		err = writer.WriteMessages(context.Background(), msg)
		if err != nil {
			log.Printf("failed to write message: %v", err)
			continue
		}

		log.Printf("sent otder %s", order.OrderUID)
		time.Sleep(1 * time.Second)
	}
}

func generateOrder() *entities.Order {
	return &entities.Order{
		OrderUID:        uuid.NewString(),
		TrackNumber:     "WBILM" + randomString(6),
		Entry:           "WBIL",
		Locale:          "en",
		CustomerID:      "test",
		DeliveryService: "meest",
		DateCreated:     time.Now().UTC(),
		Delivery: entities.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2534232432",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: entities.Payment{
			Transaction:  uuid.NewString(),
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       rand.Intn(5000) + 100,
			PaymentDt:    int(time.Now().Unix()),
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   500,
			CustomFee:    0,
		},
		Items: []entities.Item{
			{
				ChrtID:      rand.Int(),
				TrackNumber: "WBILMTESTTRACK",
				Price:       500,
				Rid:         uuid.NewString(),
				Name:        "Random Product",
				Sale:        10,
				Size:        "M",
				TotalPrice:  450,
				NmID:        rand.Int(),
				Brand:       "BrandX",
				Status:      202,
			},
		},
	}
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
