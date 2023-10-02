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
		Brokers:           "broker:9092",
		GroupId:           "Poc.Kafka.Go",
		EnableIdempotence: true,
		Offset:            enums.Earliest,
		SaslMechanism:     enums.Plain,
		SecurityProtocol:  enums.SaslSsl,
		UserName:          "ApiKey",
		Password:          "Secret",
		RequestTimeout:    5000,
		SchemaRegistry: options.SchemaRegistryOptions{
			Url:                        "https://psrc-dz0xz.westus2.azure.confluent.cloud",
			BasicAuthUser:              "SchemaRegistryApiKey",
			BasicAuthSecret:            "SchemaRegitrySecret",
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
