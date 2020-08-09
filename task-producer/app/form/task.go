package form

import "time"

type TaskForm struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}
