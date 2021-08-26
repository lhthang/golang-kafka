package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"log"
	"os"
	"os/signal"
	"task-consumer/app/message"
	"task-consumer/app/model"
	"task-consumer/app/repository"
	"time"
)

type KafkaConsumer struct {
	Consumer *cluster.Consumer
}

var Consumer KafkaConsumer

var brokers = []string{"127.0.0.1:9092", "127.0.0.1:9093"}

const (
	topic = "tasks"
)

func InitConsumer() (KafkaConsumer, error) {
	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()

	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.CommitInterval = time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// init consumer
	topics := []string{topic}
	consumer, err := cluster.NewConsumer(brokers, "group-task", topics, config)

	if err != nil {
		return KafkaConsumer{}, err
	}

	Consumer = KafkaConsumer{Consumer: consumer}
	return Consumer, nil
}

func (ks KafkaConsumer) Consume() {
	defer ks.Consumer.Close()
	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range ks.Consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume messages, watch signals
	for {
		select {
		case msg, ok := <-ks.Consumer.Messages():
			if ok {
				content:=""
				if msg.Topic != topic {
					continue
				}
				//handle msg here
				//custom or refactor code to be cleaner
				var requestMsg message.RequestMsg
				err := json.Unmarshal(msg.Value, &requestMsg)

				if err != nil {
					log.Println(err)
					continue
				}
				switch requestMsg.Method {
				case "POST":
					var task model.Task
					err = json.Unmarshal([]byte(requestMsg.Message), &task)
					if err != nil {
						break
					}
					content =repository.TaskEntity.CreateTask(requestMsg.Id, task)

				case "DELETE":
				case "UPDATE":
				case "GET":
				case "GETDETAIL":
				}

				log.Println(content)
				err = Producer.Publish(content)
				if err != nil {
					log.Println(err)
					continue
				}
				ks.Consumer.MarkOffset(msg, "") // mark message as processed
				ks.Consumer.CommitOffsets()
			}
		case <-signals:
			return
		}
	}
}
