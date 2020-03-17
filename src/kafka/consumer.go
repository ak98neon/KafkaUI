package kafka

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

var ResponseTopic = "routing-response-ofsc-5dbb3b.test"

type Info struct {
	Brokers         []string `json:"brokers"`
	Topics          []string `json:"topics"`
	CountTopics     int      `json:"count_topics"`
	CountPartitions int32    `json:"count_partitions"`
}

func ConsumeMessage() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Create new consumer
	master, err := sarama.NewConsumer(BrokerList, config)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	// How to decide partition, is it fixed value...?
	consumer, err := master.ConsumePartition(ResponseTopic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Count how many message processed
	msgCount := 0

	// Get signnal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Println("Received messages", string(msg.Key), string(msg.Value))
			case <-signals:
				fmt.Println("Interrupt is detect")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCount, "messages")
}

func CreateInfoConsumer() (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	return sarama.NewConsumer(BrokerList, config)
}

func GetKafkaInfo() Info {
	consumer, _ := CreateInfoConsumer()

	defer func() {
		if err := consumer.Close(); err != nil {
			panic(err)
		}
	}()

	topics, err := consumer.Topics()
	if err != nil {
		log.Println(err)
	}

	countTopics := len(topics)

	var countPartitions int32 = 0
	for i := 0; i < countTopics; i++ {
		partitions, _ := consumer.Partitions(topics[i])
		for i := 0; i < len(partitions); i++ {
			countPartitions += partitions[i]
		}
	}

	return Info{
		Brokers:         BrokerList,
		Topics:          topics,
		CountTopics:     countTopics,
		CountPartitions: countPartitions,
	}
}
