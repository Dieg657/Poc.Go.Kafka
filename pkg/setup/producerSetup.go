package setup

import (
	"strconv"

	"github.com/Dieg657/Poc.Go.Kafka/pkg/options"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type IProducerSetup interface {
	New(configs *options.KafkaOptions) error
	GetProducer() *kafka.Producer
}

type ProducerSetup struct {
	producerKafka *kafka.Producer
}

func (producerSetup *ProducerSetup) New(configs *options.KafkaOptions) error {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":                     configs.Brokers,
		"enable.idempotence":                    strconv.FormatBool(configs.EnableIdempotence),
		"max.in.flight.requests.per.connection": "5",
		"acks":                                  "all",
		"client.id":                             configs.GroupId,
		"sasl.mechanisms":                       string(configs.SaslMechanism),
		"security.protocol":                     string(configs.SecurityProtocol),
		"sasl.username":                         configs.UserName,
		"sasl.password":                         configs.Password,
		"request.timeout.ms":                    strconv.FormatInt(int64(configs.RequestTimeout*10), 10),
	})

	if err != nil {
		return err
	}

	producerSetup.producerKafka = producer
	return nil
}

func (producerSetup *ProducerSetup) GetProducer() *kafka.Producer {
	return producerSetup.producerKafka
}
