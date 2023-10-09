package ioc

import (
	"fmt"
	"os"

	"github.com/Dieg657/Poc.Go.Kafka/pkg/enums"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/options"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/setup"
)

type IKafkaIoC interface {
	New() error
	GetConsumer() (setup.IConsumerSetup, setup.ISchemaRegistrySetup)
	GetProducer() (setup.IProducerSetup, setup.ISchemaRegistrySetup)
}

type KafkaIoC struct {
	consumerSetup       *setup.ConsumerSetup
	producerSetup       *setup.ProducerSetup
	schemaRegistrySetup *setup.SchemaRegistrySetup
}

func (ioc *KafkaIoC) New() error {
	optionsKafka := &options.KafkaOptions{
		Brokers:           os.Getenv("KAFKA_BROKERS"),
		GroupId:           os.Getenv("KAFKA_GROUPID"),
		EnableIdempotence: true,
		Offset:            enums.AutoOffsetReset(os.Getenv("KAFKA_AUTO_OFFSET_RESET")),
		SaslMechanism:     enums.Plain,
		SecurityProtocol:  enums.SaslSsl,
		UserName:          os.Getenv("KAFKA_USERNAME"),
		Password:          os.Getenv("KAFKA_PASSWORD"),
		RequestTimeout:    5000,
		SchemaRegistry: options.SchemaRegistryOptions{
			Url:                        os.Getenv("KAFKA_SCHEMA_REGISTRY_URL"),
			BasicAuthUser:              os.Getenv("KAFKA_SCHEMA_REGISTRY_USERNAME"),
			BasicAuthSecret:            os.Getenv("KAFKA_SCHEMA_REGISTRY_PASSWORD"),
			AutoRegisterSchemas:        false,
			RequestTimeout:             5000,
			BasicAuthCredentialsSource: enums.UserInfo,
		},
	}

	schemaRegistrySetup := setup.SchemaRegistrySetup{}
	err := schemaRegistrySetup.New(optionsKafka)
	if err != nil {
		fmt.Println("Error on initialize Schema Registry")
		os.Exit(-1)
	}

	producerSetup := setup.ProducerSetup{}
	err = producerSetup.New(optionsKafka)

	if err != nil {
		fmt.Println("Eror")
		os.Exit(-1)
	}

	consumerSetup := setup.ConsumerSetup{}
	err = consumerSetup.New(optionsKafka)

	if err != nil {
		fmt.Println("Error on initialize Schema Registry")
		os.Exit(-1)
	}

	ioc.consumerSetup = &consumerSetup
	ioc.producerSetup = &producerSetup
	ioc.schemaRegistrySetup = &schemaRegistrySetup
	return nil
}

func (ioc *KafkaIoC) GetConsumer() (setup.IConsumerSetup, setup.ISchemaRegistrySetup) {
	return ioc.consumerSetup, ioc.schemaRegistrySetup
}

func (ioc *KafkaIoC) GetProducer() (setup.IProducerSetup, setup.ISchemaRegistrySetup) {
	return ioc.producerSetup, ioc.schemaRegistrySetup
}
