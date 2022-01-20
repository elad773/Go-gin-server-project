package server

import (
	. "backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	p := r.Group("/api/people/")
	{
		p.GET("", HandleGetPeople)
		p.POST("", HandlePostPerson)
		id := p.Group("/:id")
		{
			id.GET("", HandleGetPerson)
			id.PATCH("", HandlePatchPerson)
			id.DELETE("", HandleDeletePerson)
			id.GET("/tasks/", HandleGetTasks)
			id.POST("/tasks/", HandlePostTask)
		}
	}
	t := r.Group("/api/tasks/:id")
	{
		t.GET("", HandleGetTask)
		t.PATCH("", HandlePatchTask)
		t.DELETE("", HandleDeleteTask)
		t.GET("/status", HandleGetTaskStatus)
		t.PUT("/status", HandlePutTaskStatus)
		t.GET("/owner", HandleGetTaskOwnerId)
		t.PUT("/owner", HandlePutTaskOwnerId)
	}
}
