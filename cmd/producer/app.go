package producer

import (
	"context"
	"fmt"
	"sync"
	"time"

	avro "github.com/Dieg657/Poc.Go.Kafka/cmd/producer/pkg/messages"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/base"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/enums"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/ioc"
	publisher_kafka "github.com/Dieg657/Poc.Go.Kafka/pkg/publisher"
)

type Producer struct{}

func (worker *Producer) StartProducer(iocKafka ioc.IKafkaIoC, ctx context.Context, wg *sync.WaitGroup) {
	publisher := publisher_kafka.TopicPublisher[avro.MensagemSimples]{}
	publisher.New(iocKafka.GetProducer())

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Terminate signal detected... Stopping producer worker")
			wg.Done()
			return
		default:
			fmt.Printf("Enviando a mensagem\n")
			payload := base.Message[avro.MensagemSimples]{}
			payload.New(&avro.MensagemSimples{Hora_mensagem: time.Now().Format("15:04:05 02/01/2006")})
			err := publisher.Publish("teste2", payload, enums.AvroSerialization)

			if err != nil {
				fmt.Println(err)
			}

			time.Sleep(150 * time.Millisecond)
		}
	}
}
