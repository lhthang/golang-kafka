package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"task-producer/app/form"
	message2 "task-producer/app/message"
	"task-producer/app/model"
	"task-producer/kafka"
	"task-producer/utils/serialize"
	"task-producer/utils/uuid"
	"time"
)

var TaskEntity ITask

type ITask interface {
	GetAll() error
	CreateTask(method string, form form.TaskForm) (id string, err error)
}

type taskEntity struct {
}

func NewTaskEntity() ITask {
	TaskEntity = taskEntity{}
	return TaskEntity
}

func (entity taskEntity) GetAll() error {
	message := message2.RequestMsg{
		Method:  "GET",
		Message: "",
	}

	content, err := serialize.Serialize(message)

	if err != nil {
		return err
	}
	err = kafka.Producer.Publish(content)
	if err != nil {
		return err
	}

	return nil
}

func (entity taskEntity) CreateTask(method string, taskForm form.TaskForm) (string, error) {

	task := model.Task{
		Id:        primitive.NewObjectID(),
		Name:      taskForm.Name,
		CreatedAt: time.Now(),
	}

	id := uuid.GenerateUUID()
	sendMsg := message2.RequestMsg{
		Method: method,
		Id:     id,
	}

	content, err := serialize.Serialize(task)

	if err != nil {
		return id, err
	}
	sendMsg.Message = content
	content, err = serialize.Serialize(sendMsg)

	if err != nil {
		return id, err
	}

	err = kafka.Producer.Publish(content)
	if err != nil {
		return id, err
	}
	return id, nil
}
