package eventsender

import (
	configparser "companies/cmd/internal/configParser"
	"companies/cmd/internal/consts"
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type EventSender interface {
	PublishEvent(string, Event) error
}

type sender struct {
	producer *kafka.Producer
}

func NewEventSender(config configparser.Kafka) EventSender {
	log.Println(consts.ApplicationPrefix, "Starting EventSender")

	brokerAddr := configparser.GetCfgValue("KAFKA_BROKER", config.Broker)

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokerAddr})
	if err != nil {
		log.Println(consts.ApplicationPrefix, "Failed to create producer: ", err)
		return nil
	}

	s := sender{p}

	s.waitRediness()

	return &s
}

func (s *sender) waitRediness() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Println(consts.ApplicationPrefix, "Timeout waiting for Kafka to be ready")
			return
		default:
			_, err := s.producer.GetMetadata(nil, false, 1000)
			if err == nil {
				log.Println(consts.ApplicationPrefix, "Kafka is ready")
				return
			}
			log.Println(consts.ApplicationPrefix, "Waiting for Kafka broker...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (s *sender) PublishEvent(topic string, event Event) error {
	deliveryChan := make(chan kafka.Event)

	message, err := json.Marshal(event)

	err = s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	if err != nil {
		log.Fatalf("Produce error: %s", err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Println(consts.ApplicationPrefix, "Delivery failed: ", m.TopicPartition.Error)
		return errors.New("Delivery failed")
	} else {
		log.Println(consts.ApplicationPrefix, "Message delivered to ", m.TopicPartition)
	}

	close(deliveryChan)

	return nil
}

func (s *sender) Close() error {
	s.producer.Close()
	return nil
}
