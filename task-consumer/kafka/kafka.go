package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"log"
	"os"
	"os/signal"
	"task-consumer/app/repository"
	"time"
)

var brokers = []string{"127.0.0.1:9092", "127.0.0.1:9093"}

const (
	topic = "tasks"
)

func InitConsumer() (*cluster.Consumer, error) {
	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()

	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.CommitInterval = time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// init consumer
	topics := []string{topic}
	consumer, err := cluster.NewConsumer(brokers, "group-task", topics, config)

	if err != nil {
		return nil, err
	}

	return consumer, err
}

func Consume(consumer *cluster.Consumer) {
	defer consumer.Close()
	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume messages, watch signals
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				if msg.Topic != topic {
					continue
				}
				//handle msg here
				//custom or refactor code to be cleaner
				var postMsg MethodMsg
				err := json.Unmarshal(msg.Value, &postMsg)

				if err != nil {
					log.Println(err)
					continue
				}
				switch postMsg.Method {
				case "POST":
					var content PostMsg
					err = json.Unmarshal(msg.Value, &content)
					if err != nil {
						break
					}
					log.Println(postMsg.Method, content.Task)
					repository.TaskEntity.CreateTask(content.Task)
				case "DELETE":
				case "UPDATE":
				}
				consumer.MarkOffset(msg, "") // mark message as processed
			}
		case <-signals:
			return
		}
	}
}
