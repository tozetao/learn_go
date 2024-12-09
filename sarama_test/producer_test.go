package sarama_test

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"os/signal"
	"sync"
	"testing"
)

var addr = []string{"localhost:9094"}

// 测试同步的生产者
func TestSyncProducer(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(addr, config)
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	require.NoError(t, err)

	defer func() {
		if err := producer.Close(); err != nil {
			t.Logf("producer close error, %v", err)
		}
	}()

	for i := 0; i < 3; i++ {
		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: "my_topic",
			Value: sarama.StringEncoder("hello."),
		})
		assert.NoError(t, err)
		t.Logf("partition = %d, offset = %d", partition, offset)
	}

}

// 测试异步的生产者
func TestAsyncProducer(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(addr, config)
	if err != nil {
		panic(err)
	}

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var (
		wg                                  sync.WaitGroup
		enqueued, successes, producerErrors int
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			log.Println(err)
			producerErrors++
		}
	}()

ProducerLoop:
	for {
		message := &sarama.ProducerMessage{Topic: "my_topic", Value: sarama.StringEncoder("testing 123")}
		select {
		case producer.Input() <- message:
			enqueued++

		case <-signals:
			producer.AsyncClose() // Trigger a shutdown of the producer.
			break ProducerLoop
		}
	}

	wg.Wait()

	log.Printf("Successfully produced: %d; errors: %d\n", successes, producerErrors)
}

func TestAsyncProducerV1(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(addr, config)
	require.NoError(t, err)

	input := producer.Input()
	input <- &sarama.ProducerMessage{
		Topic: "my_topic",
		Value: sarama.StringEncoder("hello world"),
	}

	select {
	case err = <-producer.Errors():
		t.Logf("deviced error, %v", err)
	case success := <-producer.Successes():
		t.Logf("deviced success, %v", success.Value)
	}
}
