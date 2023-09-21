package publisher

import (
	"errors"
	"fmt"

	"github.com/Dieg657/Poc.Go.Kafka/pkg/base"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/enums"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/setup"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ITopicPublisher[TData any] interface {
	New(producerSetup setup.IProducerSetup, registry *setup.ISchemaRegistrySetup) error
	Publish(topic *string, message *base.Message[TData], serialization enums.MessageSerialization) error
}

type TopicPublisher[TData any] struct {
	producer *kafka.Producer
	registry setup.ISchemaRegistrySetup
}

func (publisher *TopicPublisher[TData]) New(producerSetup setup.IProducerSetup, registry setup.ISchemaRegistrySetup) error {
	if producerSetup == nil {
		return errors.New("There's no producer configurated!")
	}

	if registry == nil {
		return errors.New("There's no schema registry configured!")
	}

	publisher.producer = producerSetup.GetProducer()
	publisher.registry = registry

	return nil
}

func (publisher *TopicPublisher[TData]) Publish(topic string, message base.Message[TData], serialization enums.MessageSerialization) error {

	payload, err := publisher.serializePayload(topic, *message.GetData(), serialization)

	if err != nil {
		fmt.Println("Failed attempt to serialize message")
		return err
	}

	err = publisher.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
		Headers:        []kafka.Header{{Key: "correlationId", Value: []byte((message.Header.CorrelationId).String())}},
		Key:            []byte((message.Header.CorrelationId).String()),
	}, nil)

	if err != nil {
		fmt.Printf("Failed when produce message: %v\n", err)
		return err
	}

	publisher.producer.Flush(10)
	return nil
}

func (publisher *TopicPublisher[TData]) serializePayload(topic string, payload TData, serialization enums.MessageSerialization) ([]byte, error) {

	switch serialization {
	case enums.AvroSerialization:
		return publisher.registry.GetAvroSerializer().Serialize(topic, &payload)
	case enums.JsonApiSerialization:
	case enums.JsonStandardSerialization:
		return publisher.registry.GetJsonSerializer().Serialize(topic, &payload)
	default:
		return nil, errors.New("Invalid serialization type")
	}

	return nil, errors.New("Invalid operation exception")
}
