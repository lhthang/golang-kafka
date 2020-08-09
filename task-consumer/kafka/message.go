package kafka

import "task-consumer/app/model"

type MethodMsg struct {
	Method string `json:"method"`
}

type PostMsg struct {
	Task model.Task `json:"task"`
}

type DeleteMsg struct {
	Id string `json:"id"`
}

type UpdateMsg struct {
	//define here
}
