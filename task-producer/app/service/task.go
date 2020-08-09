package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"task-producer/app/form"
	"task-producer/app/model"
	"task-producer/kafka"
	"task-producer/utils/serialize"
	"time"
)

var TaskEntity ITask

type ITask interface {
	GetAll() ([]form.TaskForm, int, error)
	CreateTask(method string, form form.TaskForm) (*model.Task, int, error)
}

type taskEntity struct {
}

func NewTaskEntity() ITask {
	TaskEntity = taskEntity{}
	return TaskEntity
}

func (entity taskEntity) GetAll() ([]form.TaskForm, int, error) {
	message := kafka.Message{
		Method: "GET",
		Id:     "",
	}

	content, err := serialize.Serialize(message)

	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	err = kafka.Producer.Publish(content)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nil, http.StatusOK, nil
}

func (entity taskEntity) CreateTask(method string, taskForm form.TaskForm) (*model.Task, int, error) {

	task := model.Task{
		Id:        primitive.NewObjectID(),
		Name:      taskForm.Name,
		CreatedAt: time.Now(),
	}

	message := kafka.Message{
		Method: method,
		Id:     "",
		Task:   &task,
	}

	content, err := serialize.Serialize(message)

	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	err = kafka.Producer.Publish(content)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	return &task, http.StatusOK, nil
}
