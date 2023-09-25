package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	avro "github.com/Dieg657/Poc.Go.Kafka/cmd/consumer/pkg/messages"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/base"
	consumer_kafka "github.com/Dieg657/Poc.Go.Kafka/pkg/consumer"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/enums"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/ioc"
)

type Consumer struct{}

func (worker *Consumer) StartConsumer(kafkaIoc ioc.IKafkaIoC, ctx context.Context, wg *sync.WaitGroup) {
	consumer := consumer_kafka.TopicConsumer[avro.MensagemSimples]{}
	consumer.New(kafkaIoc.GetConsumer())
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping consumer worker...")
			wg.Done()
			return
		default:
			err := consumer.Consume("teste2", enums.AvroDeserialization, func(message base.Message[avro.MensagemSimples]) error {
				fmt.Printf("Foi recebido o correlationId: %s", message.CorrelationId)
				fmt.Printf("Hora da mensagem recebida: %s\n", message.Data.Hora_mensagem)
				time.Sleep(5 * time.Millisecond)
				return nil
			})

			if err != nil {
				fmt.Println("Erro ao consumir mensagem!")
			}
		}
	}
}
