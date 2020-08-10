package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
)

const (
	kafkaConn = "localhost:9092"
	sendToTopic = "r-tasks"
)


type KafkaProducer struct {
	Producer sarama.SyncProducer
}

var Producer KafkaProducer

func InitProducer()(KafkaProducer, error) {
	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.ClientID="group1"

	// async producer
	//prd, err := sarama.NewAsyncProducer([]string{kafkaConn}, config)

	// sync producer
	prd, err := sarama.NewSyncProducer([]string{kafkaConn}, config)

	Producer = KafkaProducer{Producer: prd}
	return Producer, err
}

func(kafka KafkaProducer) Publish(message string) error {
	// publish sync
	msg := &sarama.ProducerMessage {
		Topic: sendToTopic,
		Value: sarama.StringEncoder(message),
	}
	p, o, err := kafka.Producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error publish: ", err.Error())
		return err
	}

	// publish async
	//producer.Input() <- &sarama.ProducerMessage{

	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)
	return nil
}

