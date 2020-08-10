package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"task-producer/app/api"
	"task-producer/kafka"
	"task-producer/middlewares"
)

type Server struct {
}

func (s Server) StartServer() {
	r := gin.Default()

	//Apply CORs before use Group
	r.Use(middlewares.NewCors([]string{"*"}))
	r.Use(gin.Logger())
	r.Use(middlewares.NewRecovery())

	publicRoute := r.Group("/api/v1")

	_, err := kafka.InitProducer()
	if err != nil {
		logrus.Println(err)
		os.Exit(1)
	}
	// init consumer
	_, err = kafka.InitConsumer()
	if err != nil {
		logrus.Println("Error consumer group: ", err.Error())
	}
	defer kafka.Consumer.Consumer.Close()
	api.ApplyTaskAPI(publicRoute)

	r.Run(":8989")
}
