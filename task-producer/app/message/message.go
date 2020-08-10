package message

import "task-producer/app/model"

type RequestMsg struct {
	Method  string `json:"method"`
	Id      string `json:"id"`
	Message string `json:"message"`
}

type ReceiveMsg struct {
	Type    string `json:"type"`
	Err     error  `json:"error"`
	Code    int    `json:"code"`
	Id      string `json:"id"`
	Message string `json:"message"`
}

type ReceiveObjMsg struct {
	Task model.Task `json:"task"`
}

type ReceiveArrMsg struct {
	Tasks []model.Task `json:"tasks"`
}
