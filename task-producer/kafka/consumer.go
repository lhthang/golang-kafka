package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"log"
	"os"
	"os/signal"
	"task-producer/app/message"
	"task-producer/app/model"
	"time"
)

type KafkaConsumer struct {
	Consumer *cluster.Consumer
}

var Consumer KafkaConsumer

var brokers = []string{"127.0.0.1:9092", "127.0.0.1:9093"}

const (
	receiveTopic = "r-tasks"
)

func InitConsumer() (KafkaConsumer, error) {
	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()

	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.CommitInterval = time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// init consumer
	topics := []string{receiveTopic}
	consumer, err := cluster.NewConsumer(brokers, "group-task", topics, config)

	if err != nil {
		return KafkaConsumer{}, err
	}

	Consumer = KafkaConsumer{Consumer: consumer}
	return Consumer, nil
}

func (ks KafkaConsumer) Consume(id string) (v interface{}, code int, err error) {
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
			code := 200
			var v interface{}
			var errResp error
			log.Println(msg.Topic)
			if ok {
				if msg.Topic != receiveTopic {
					continue
				}
				//handle msg here
				//custom or refactor code to be cleaner
				var recvMsg message.ReceiveMsg
				err := json.Unmarshal(msg.Value, &recvMsg)

				if err != nil {
					log.Println(err)
					continue
				}
				if recvMsg.Id != id {
					continue
				}
				code = recvMsg.Code
				errResp = recvMsg.Err

				switch recvMsg.Type {
				case "Object":
					var content model.Task
					err = json.Unmarshal([]byte(recvMsg.Message), &content)
					if err != nil {
						continue
					}
					v = content
				case "Array":
					var content []model.Task
					err = json.Unmarshal([]byte(recvMsg.Message), &content)
					if err != nil {
						continue
					}
					v = content
				}
				log.Println(v)
				ks.Consumer.MarkOffset(msg, "") // mark message as processed
				return v, code, errResp
			}
		case <-signals:
			return
		}
	}
}
