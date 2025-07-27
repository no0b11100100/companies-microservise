package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/require"
)

func TestCreateCompanyIntegration(t *testing.T) {
	// Build request body
	reqBody := map[string]interface{}{
		"name":         "Test Company",
		"employees":    42,
		"isRegistered": true,
		"type":         2,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/companies", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	token, err := GenerateToken("root")
	if err != nil {
		t.Log("Can not create token")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token)) // add if auth middleware enabled
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify Kafka event published

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBootstrap,
		"group.id":          "test-group",
		"auto.offset.reset": "earliest",
	})
	require.NoError(t, err)
	defer consumer.Close()

	err = consumer.Subscribe("data-changed", nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msgChan := make(chan *kafka.Message, 1)
	go func() {
		for {
			ev := consumer.Poll(500)
			switch e := ev.(type) {
			case *kafka.Message:
				msgChan <- e
				return
			}
		}
	}()

	select {
	case msg := <-msgChan:
		require.Contains(t, string(msg.Value), "Test Company")
	case <-ctx.Done():
		t.Fatal("Kafka message was not received in time")
	}
}
