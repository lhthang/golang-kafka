package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"task-producer/app/form"
	"task-producer/app/service"
	"task-producer/kafka"
)

func ApplyTaskAPI(r *gin.RouterGroup) {
	taskEntity := service.NewTaskEntity()
	taskRoute := r.Group("/tasks")
	taskRoute.GET("", getAllTasks())
	taskRoute.POST("", createTask(taskEntity))
}

func getAllTasks() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
	}
}

func createTask(entity service.ITask) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var taskForm form.TaskForm

		if err := ctx.Bind(&taskForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		id,err := entity.CreateTask(ctx.Request.Method, taskForm)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		var resp interface{}
		var code int

		resp, code, err = kafka.Consumer.Consume(id)

		response := map[string]interface{}{
			"task":  resp,
			"error": err,
		}
		ctx.JSON(code, response)
	}
}


