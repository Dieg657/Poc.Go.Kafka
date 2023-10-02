package setup

import (
	"os"

	"github.com/Dieg657/Poc.Go.Kafka/pkg/options"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type IConsumerSetup interface {
	New(configs *options.KafkaOptions) error
	GetConsumer() *kafka.Consumer
}

type ConsumerSetup struct {
	consumerKafka *kafka.Consumer
}

func (consumerSetup *ConsumerSetup) New(configs *options.KafkaOptions) error {
	hostname, err := os.Hostname()

	if err != nil {
		return err
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  configs.Brokers,
		"group.id":           configs.GroupId,
		"client.id":          hostname,
		"sasl.mechanisms":    string(configs.SaslMechanism),
		"security.protocol":  string(configs.SecurityProtocol),
		"sasl.username":      configs.UserName,
		"sasl.password":      configs.Password,
		"session.timeout.ms": 45000,
		"enable.auto.commit": "true",
		"auto.offset.reset":  string(configs.Offset),
	})

	if err != nil {
		return err
	}

	consumerSetup.consumerKafka = consumer
	return nil
}

func (consumerSetup *ConsumerSetup) GetConsumer() *kafka.Consumer {
	return consumerSetup.consumerKafka
}
