package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

type TokenResponse struct {
	Token string `json:"token"`
}

func getToken(appURL string) string {
	timeout := time.After(60 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("Timeout reached. Token not received.")
			return ""
		case <-ticker.C:
			req, err := http.NewRequest("POST", appURL+"/api/v1/token", nil)
			if err != nil {
				log.Println("Failed to create HTTP request:", err)
				continue
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Println("HTTP POST request failed:", err)
				continue
			}

			defer resp.Body.Close()

			var tokenData TokenResponse
			if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
				log.Println("Failed to parse JSON:", err)
				continue
			}

			if tokenData.Token != "" {
				log.Println("Token received:", tokenData.Token)
				return tokenData.Token
			}

			log.Println("Waiting for token...")
		}
	}
}

func TestIntegration_CreateRecord(t *testing.T) {
	appHost := getEnv("APP_HOST", "app")
	appPort := getEnv("APP_PORT", "8080")
	appURL := fmt.Sprintf("http://%s:%s", appHost, appPort)

	kafkaBroker := getEnv("KAFKA_BROKER", "kafka:9092")
	kafkaTopic := "data-changed"

	t.Log("Waiting for services to be ready...")
	jwtToken := getToken(appURL)

	if len(jwtToken) == 0 {
		t.Fatal("Can not reach service")
	}

	requestURL := appURL + "/api/v1/companies/"
	t.Logf("Testing HTTP POST to %s", requestURL)
	postBody := map[string]interface{}{
		"name":           "Cool Company",
		"employeesCount": 42,
		"isRegistered":   true,
		"type":           1,
	}

	jsonBytes, err := json.Marshal(postBody)
	if err != nil {
		t.Fatalf("failed to marshal POST body: %v", err)
	}

	req, err := http.NewRequest("POST", appURL+"/api/v1/companies/", bytes.NewReader(jsonBytes))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("HTTP POST request failed: %v", err)
	}

	defer resp.Body.Close()

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
