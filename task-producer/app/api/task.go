package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"task-producer/app/form"
	"task-producer/app/service"
	error2 "task-producer/utils/error"
)

func ApplyTaskAPI(r *gin.RouterGroup) {
	taskEntity :=service.NewTaskEntity()
	taskRoute := r.Group("/tasks")
	taskRoute.GET("", getAllTasks())
	taskRoute.POST("", createTask(taskEntity))
}

func getAllTasks() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		response := map[string]interface{}{
			"tasks": nil,
		}
		ctx.JSON(200, response)
	}
}

func createTask(entity service.ITask) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var taskForm form.TaskForm

		if err := ctx.Bind(&taskForm); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		task,code,err:=entity.CreateTask(ctx.Request.Method,taskForm)
		response := map[string]interface{}{
			"task": task,
			"error":error2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}