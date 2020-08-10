package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"task-consumer/app/message"
	"task-consumer/app/model"
	"task-consumer/my_db"
	"task-consumer/utils/serialize"
)

var TaskEntity ITask

type taskEntity struct {
	resource *my_db.Resource
	repo     *mongo.Collection
}

type ITask interface {
	CreateTask(id string, task model.Task) string
}

func NewTaskEntity(rs *my_db.Resource) ITask {

	TaskEntity = &taskEntity{
		resource: rs,
		repo:     rs.DB.Collection("tasks"),
	}
	return TaskEntity
}

func (entity taskEntity) CreateTask(id string, task model.Task) string {
	ctx, cancel := initContext()
	defer cancel()

	_, err := entity.repo.InsertOne(ctx, task)

	respMsg := message.ResponseMsg{
		Type:    "Object",
		Code:    0,
		Err:     nil,
		Id:      id,
		Message: "",
	}
	if err != nil {
		respMsg.Code = http.StatusBadRequest
		respMsg.Err = err
	}

	content, err := serialize.Serialize(task)

	if err != nil {
		log.Println(err)
	}
	respMsg.Message = content
	respMsg.Code = http.StatusOK

	content, err = serialize.Serialize(respMsg)

	if err != nil {
		log.Println(err)
	}
	return content
}
