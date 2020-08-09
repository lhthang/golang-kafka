package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"task-consumer/app/model"
	"task-consumer/my_db"
)

var TaskEntity ITask

type taskEntity struct {
	resource *my_db.Resource
	repo     *mongo.Collection
}

type ITask interface {
	CreateTask(task model.Task) (model.Task,error)
}

func NewTaskEntity(rs *my_db.Resource) ITask  {
	TaskEntity=&taskEntity{
		resource: rs,
		repo:     rs.DB.Collection("tasks"),
	}
	return TaskEntity
}

func(entity taskEntity) CreateTask(task model.Task) (model.Task,error)  {
	ctx,cancel:=initContext()
	defer cancel()

	_,err:=entity.repo.InsertOne(ctx,task)
	if err!=nil{
		return model.Task{},err
	}
	return task,nil
}
