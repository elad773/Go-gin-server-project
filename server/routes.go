package server

// func (s *Instance) routes() {
// 	s.Engine.GET("/api/people/", s.handlePeopleGet())
// 	s.Engine.POST("/api/people/", s.handlePeoplePatch())
// 	s.Engine.GET("/api/people/:id", s.handleGetPerson())
// 	s.Engine.PATCH("/api/people/:id", s.handlePersonPatch())
// 	s.Engine.DELETE("/api/people/:id", s.handlePersonDelete())
// 	s.Engine.GET("/api/people/:id/tasks/*status", s.handleTasksGet())
// 	s.Engine.POST("/api/people/:id/tasks/", s.handleTaskPost())
// 	s.Engine.GET("/api/tasks/:id", s.handleTaskGet())
// 	s.Engine.PATCH("/api/tasks/:id", s.handleTaskPatch())
// 	s.Engine.DELETE("/api/tasks/:id", s.handleTaskDelete())
// 	s.Engine.GET("/api/tasks/:id/status", s.handleTaskStatusGet())
// 	s.Engine.PUT("/api/tasks/:id/status", s.handleTaskStatusPut())
// 	s.Engine.GET("/api/tasks/:id/owner", s.handleTaskOwnerIdGet())
// 	s.Engine.PUT("/api/tasks/:id/owner", s.handleTaskOwnerIdPut())
// }

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	p := r.Group("/api/people/")
	{
		p.GET("", handleGetPeople)
		p.POST("", handlePostPerson)
		id := p.Group("/:id")
		{
			id.GET("", handleGetPerson)
			id.PATCH("", handlePatchPerson)
			id.DELETE("", handleDeletePerson)
			id.GET("/tasks/", handleGetTasks)
			id.POST("/tasks/", handlePostTask)
		}
	}
	t := r.Group("/api/tasks/:id")
	{
		t.GET("", handleGetTask)
		t.PATCH("", handlePatchTask)
		t.DELETE("", handleDeleteTask)
		t.GET("/status", handleGetTaskStatus)
		t.PUT("/status", handlePutTaskStatus)
		t.GET("/owner", handleGetTaskOwnerId)
		t.PUT("/owner", handlePutTaskOwnerId)
	}
}
