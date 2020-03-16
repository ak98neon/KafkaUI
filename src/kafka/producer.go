package kafka

import (
	"bytes"
	"fmt"
	"github.com/Shopify/sarama"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"strconv"
	"time"
)

var ProducerConfig = sarama.NewConfig()
var BrokerList []string
var ClientId string

func ProduceMessage(file multipart.File, count int, topic string) {
	ProducerConfig.Version = sarama.V0_11_0_2
	ProducerConfig.Producer.RequiredAcks = sarama.WaitForAll
	ProducerConfig.Producer.Retry.Max = 5
	ProducerConfig.Producer.MaxMessageBytes = 304857600
	ProducerConfig.Producer.Return.Successes = true
	ProducerConfig.ClientID = ClientId

	producer, err := sarama.NewSyncProducer(BrokerList, ProducerConfig)
	if err != nil {
		// Should not reach here
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			// Should not reach here
			panic(err)
		}
	}()

	bufFile := bytes.NewBuffer(nil)
	if _, err = io.Copy(bufFile, file); err != nil {
		log.Println(err)
	}

	fmt.Println("File size:", len(bufFile.Bytes()))
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(bufFile.Bytes()),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("ORIG-TIME"),
				Value: []byte(strconv.Itoa(int(time.Now().Unix()))),
			},
			{
				Key:   []byte("ORIG-TTL"),
				Value: []byte(strconv.Itoa(rand.Int())),
			},
			{
				Key:   []byte("TRACE-ID"),
				Value: []byte(strconv.Itoa(rand.Int())),
			},
		},
	}

	if count > 1 {
		for i := 0; i < count; i++ {
			sendMessage(msg, producer, topic)
		}
	} else {
		sendMessage(msg, producer, topic)
	}

}

func sendMessage(msg *sarama.ProducerMessage, producer sarama.SyncProducer, topic string) {
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
}
