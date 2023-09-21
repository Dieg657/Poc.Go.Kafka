package main

import (
	"fmt"
	"time"

	mensagem "github.com/Dieg657/Poc.Go.Kafka/cmd/consumer/pkg/messages"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/base"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/consumer"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/enums"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/options"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/publisher"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/setup"
)

func main() {
	optionsKafka := &options.KafkaOptions{
		Brokers:           "localhost:9092",
		GroupId:           "console-teste-kafka",
		EnableIdempotence: true,
		Offset:            enums.Earliest,
		UserName:          "",
		Password:          "",
		RequestTimeout:    5000,
		SchemaRegistry: options.SchemaRegistryOptions{
			Url:                        "http://localhost:8081",
			BasicAuthUser:              "user",
			BasicAuthSecret:            "user",
			AutoRegisterSchemas:        false,
			RequestTimeout:             5000,
			BasicAuthCredentialsSource: enums.UserInfo,
		},
	}

	schemaRegistry := setup.SchemaRegistrySetup{}
	schemaRegistry.New(optionsKafka)

	producerSetup := setup.ProducerSetup{}
	producerSetup.New(optionsKafka)

	publisher := publisher.TopicPublisher[mensagem.MensagemSimples]{}
	publisher.New(&producerSetup, &schemaRegistry)

	consumerSetup := setup.ConsumerSetup{}
	consumerSetup.New(optionsKafka)

	consumer := consumer.TopicConsumer[mensagem.MensagemSimples]{}
	consumer.New(&consumerSetup, &schemaRegistry)
	go func() {
		for true {
			fmt.Printf("Enviando a mensagem\n")
			payload := base.Message[mensagem.MensagemSimples]{}
			payload.New(&mensagem.MensagemSimples{Hora_mensagem: time.Now().Format("15:04:05 02/01/2006")})
			err := publisher.Publish("teste2", payload, enums.AvroSerialization)

			if err != nil {
				fmt.Println(err)
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	run := true
	for run {
		err := consumer.Consume("teste2", enums.AvroDeserialization, func(message base.Message[mensagem.MensagemSimples]) error {
			fmt.Printf("Hora da mensagem recebida: %s\n", message.GetData().Hora_mensagem)
			time.Sleep(1 * time.Second)
			return nil
		})

		if err != nil {
			fmt.Println("Erro ao consumir mensagem!")
		}
	}
}
