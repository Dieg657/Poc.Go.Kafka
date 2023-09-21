package consumer

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dieg657/Poc.Go.Kafka/pkg/base"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/enums"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/setup"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ITopicConsumer[TData any] interface {
	New(consumerSetup setup.IConsumerSetup, registry *setup.ISchemaRegistrySetup) error
	Consume(topic string, deserialization enums.MessageDeserialization, handler func(message base.Message[TData]) error) error
}

type TopicConsumer[TData any] struct {
	consumer *kafka.Consumer
	registry setup.ISchemaRegistrySetup
}

func (topicConsumer *TopicConsumer[TData]) New(consumerSetup setup.IConsumerSetup, registry setup.ISchemaRegistrySetup) error {
	if consumerSetup == nil {
		return errors.New("There's no consumer configurated!")
	}

	if registry == nil {
		return errors.New("There's no schema registry configured!")
	}

	topicConsumer.consumer = consumerSetup.GetConsumer()
	topicConsumer.registry = registry

	return nil
}

func (topicConsumer *TopicConsumer[TData]) Consume(topic string, deserialization enums.MessageDeserialization, handler func(message base.Message[TData]) error) error {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	err := topicConsumer.consumer.Subscribe(topic, rebalanceCallback)

	if err != nil {
		return err
	}

	run := true

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := topicConsumer.consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				baseMessage := base.Message[TData]{}
				err := deserializeMessage(topicConsumer, e, deserialization, baseMessage.GetData())
				if err != nil {
					fmt.Printf("Failed to deserialize payload: %s\n", err)
				} else {
					fmt.Printf("%% Message on %s:\n", e.TopicPartition)
					err := handler(baseMessage)
					if err != nil {
						fmt.Println("Error on handle message")
					}
				}
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	return nil
}

func deserializeMessage[TData any](topicConsumer *TopicConsumer[TData], e *kafka.Message, deserialization enums.MessageDeserialization, data *TData) error {
	switch deserialization {
	case enums.AvroDeserialization:
		return topicConsumer.registry.GetAvroDeserializer().DeserializeInto(*e.TopicPartition.Topic, e.Value, data)
	case enums.JsonApiDeserialization:
		return topicConsumer.registry.GetJsonDeserializer().DeserializeInto(*e.TopicPartition.Topic, e.Value, data)
	}
	return errors.New("Invalid deserializer")
}

func rebalanceCallback(c *kafka.Consumer, event kafka.Event) error {
	switch ev := event.(type) {
	case kafka.AssignedPartitions:
		fmt.Printf("%% %s rebalance: %d new partition(s) assigned: %v\n",
			c.GetRebalanceProtocol(), len(ev.Partitions), ev.Partitions)

		err := c.Assign(ev.Partitions)
		if err != nil {
			return err
		}

	case kafka.RevokedPartitions:
		fmt.Printf("%% %s rebalance: %d partition(s) revoked: %v\n",
			c.GetRebalanceProtocol(), len(ev.Partitions), ev.Partitions)

		if c.AssignmentLost() {
			fmt.Fprintln(os.Stderr, "Assignment lost involuntarily, commit may fail")
		}

		commitedOffsets, err := c.Commit()

		if err != nil && err.(kafka.Error).Code() != kafka.ErrNoOffset {
			fmt.Fprintf(os.Stderr, "Failed to commit offsets: %s\n", err)
			return err
		}
		fmt.Printf("%% Commited offsets to Kafka: %v\n", commitedOffsets)

	default:
		fmt.Fprintf(os.Stderr, "Unxpected event type: %v\n", event)
	}

	return nil
}
