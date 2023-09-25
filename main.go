package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Dieg657/Poc.Go.Kafka/cmd/consumer"
	"github.com/Dieg657/Poc.Go.Kafka/cmd/producer"
	"github.com/Dieg657/Poc.Go.Kafka/pkg/ioc"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	iocKafka := ioc.KafkaIoC{}
	err := iocKafka.New()

	if err != nil {
		fmt.Println("Error on initialize Kafka IoC")
	}

	workerProducer := producer.Producer{}
	workerConsumer := consumer.Consumer{}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	wg := new(sync.WaitGroup)

	wg.Add(2)

	go func() {
		workerProducer.StartProducer(&iocKafka, ctx, wg)
	}()

	go func() {
		workerConsumer.StartConsumer(&iocKafka, ctx, wg)
	}()

	wg.Wait()
}
