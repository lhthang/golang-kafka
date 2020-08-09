package kafka

import "task-producer/app/model"

type Message struct {
	Method string     `json:"method"`
	Id     string     `json:"id"`
	Task   *model.Task `json:"task"`
}
