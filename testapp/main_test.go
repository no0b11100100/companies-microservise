package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

func TestIntegration(t *testing.T) {
	// Read env vars with fallback
	appHost := getEnv("APP_HOST", "app")
	appPort := getEnv("APP_PORT", "8080")
	appURL := fmt.Sprintf("http://%s:%s", appHost, appPort)

	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "companiesdb_test")

	kafkaBroker := getEnv("KAFKA_BROKER", "kafka:9092")
	kafkaTopic := "data-changed"

	// Wait for services to become ready (adjust as needed)
	t.Log("Waiting for services to be ready...")
	time.Sleep(50 * time.Second)

	// 1. HTTP test: POST request with JSON body
	t.Logf("Testing HTTP POST to %s", appURL)
	postBody := map[string]interface{}{
		"field1": "value1",
		"field2": 42,
	}
	jsonBytes, err := json.Marshal(postBody)
	if err != nil {
		t.Fatalf("failed to marshal POST body: %v", err)
	}

	resp, err := http.Post(appURL, "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		t.Fatalf("HTTP POST request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected HTTP 200 OK, got %d", resp.StatusCode)
	}

	// 2. DB test: open connection and check data
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	t.Logf("Connecting to DB: %s", dbDSN)
	db, err := sql.Open("mysql", dbDSN)
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	var count int
	query := "SELECT COUNT(*) FROM your_table WHERE your_condition = ?" // replace with your real query
	err = db.QueryRow(query, "some_value").Scan(&count)
	if err != nil {
		t.Fatalf("DB query failed: %v", err)
	}
	if count == 0 {
		t.Fatal("expected at least one matching row in DB, found none")
	}

	// 3. Kafka test: consume one message
	msg := consumeKafkaMessage(t, kafkaBroker, kafkaTopic)
	if msg == "" {
		t.Fatal("expected kafka message, got empty")
	}
	t.Logf("Kafka message received: %s", msg)
}

func consumeKafkaMessage(t *testing.T, broker, topic string) string {
	config := &kafka.ConfigMap{
		"bootstrap.servers":        broker,
		"group.id":                 "test-group",
		"auto.offset.reset":        "earliest",
		"enable.auto.commit":       false,
		"go.events.channel.enable": true,
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		t.Fatalf("failed to create kafka consumer: %v", err)
	}
	defer consumer.Close()

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		t.Fatalf("failed to subscribe to topic: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			t.Fatal("timeout waiting for kafka message")
		case ev := <-consumer.Events():
			switch e := ev.(type) {
			case *kafka.Message:
				return string(e.Value)
			case kafka.Error:
				t.Fatalf("kafka error: %v", e)
			}
		}
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
