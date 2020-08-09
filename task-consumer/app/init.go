package app

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"task-consumer/app/repository"
	"task-consumer/kafka"
	"task-consumer/my_db"
)

type Server struct {

}

func(s Server) StartServer()  {

	resource, err := my_db.InitResource()
	if err != nil {
		logrus.Error(err)
	}
	defer resource.Close()

	repository.NewTaskEntity(resource)

	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// init consumer
	cg, err := kafka.InitConsumer()
	if err != nil {
		logrus.Println("Error consumer group: ", err.Error())
	}

	// run consumer
	kafka.Consume(cg)
}