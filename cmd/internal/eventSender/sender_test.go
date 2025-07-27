package eventsender

import (
	"errors"
	"testing"

	"companies/cmd/internal/structs"
	"companies/cmd/tests/mocks"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var dummyEvent = structs.Event{
	URL:    "/test",
	Type:   structs.Created,
	Status: structs.Success,
}

func TestSender_PublishEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockProducer := mocks.NewMockProducer(ctrl)

	s := &sender{producer: mockProducer}

	mockProducer.EXPECT().Produce(gomock.Any(), gomock.Any()).DoAndReturn(
		func(msg *kafka.Message, deliveryChan chan kafka.Event) error {
			go func() {
				deliveryChan <- &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic:     msg.TopicPartition.Topic,
						Partition: 0,
						Error:     nil,
					},
					Value: msg.Value,
				}
			}()
			return nil
		},
	)

	err := s.PublishEvent("test-topic", dummyEvent)
	require.NoError(t, err)
}

func TestSender_PublishEvent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockProducer := mocks.NewMockProducer(ctrl)

	s := &sender{producer: mockProducer}

	mockProducer.EXPECT().Produce(gomock.Any(), gomock.Any()).DoAndReturn(
		func(msg *kafka.Message, deliveryChan chan kafka.Event) error {
			go func() {
				deliveryChan <- &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic:     msg.TopicPartition.Topic,
						Partition: 0,
						Error:     errors.New("some delivery error"),
					},
					Value: msg.Value,
				}
			}()
			return nil
		},
	)

	err := s.PublishEvent("test-topic", dummyEvent)
	require.Error(t, err)
	require.Equal(t, "Delivery failed", err.Error())
}

func TestSender_Close(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockProducer := mocks.NewMockProducer(ctrl)
	mockProducer.EXPECT().Close().Return()

	s := &sender{producer: mockProducer}
	err := s.Close()
	require.NoError(t, err)
}
