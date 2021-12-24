package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func handlePostTask(c *gin.Context) {
	id := c.Param("id")
	var person Person
	var task Task
	create := struct {
		Title   string `json:"title" binding:"required"`
		Details string `json:"details" binding:"required"`
		DueDate string `json:"dueDate" binding:"required"`
		Status  string `json:"status" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&create); err != nil {
		c.String(http.StatusBadRequest, "Required data fields are missing, data makes no sense, or data contains illegal values.")
		return
	}
	if _, err := time.Parse("2006-01-02", create.DueDate); err != nil {
		c.String(http.StatusBadRequest, "dueDate"+create.DueDate+"not in time format.")

		return
	}

	if result := person.GetPersonDB(id); result.Error != nil {

		c.String(http.StatusNotFound, "A person with the id '%s'does not exist. ", id)
		return
	}
	params := OptParams{Title: &create.Title, Details: &create.Details, DueDate: &create.DueDate, Status: create.Status}
	result := task.CreateTaskDB(&person, &params)
	if result != nil {
		c.String(http.StatusBadRequest, "%s", result.Error())
		return
	}
	c.Header("Location", fmt.Sprintf("/api/tasks/%s", task.ID))
	c.Header("x-Created-Id", fmt.Sprint(task.ID))
	c.String(http.StatusCreated, "Task created successfuly.")

}

func handleGetTasks(c *gin.Context) {

	id := c.Param("id")
	status := c.Query("status")
	var person Person
	if result := person.GetPersonDB(id); result.Error != nil {
		c.String(http.StatusNotFound, "A person with the id '%s' does not exist.", id)
		return
	}
	person.GetTasksDB(status)
	c.JSON(http.StatusOK, person.Tasks)

}

func handleGetTask(c *gin.Context) {
	id := c.Param("id")
	var task Task
	if result := task.GetTaskDB(id); result.Error != nil {
		c.String(http.StatusNotFound, "A Task with the id '%s' does not exist.", id)
		return
	}
	c.JSON(http.StatusOK, task)
}

func handlePatchTask(c *gin.Context) {
	var params OptParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(http.StatusBadRequest, "Format not valid")
		return
	}

	id := c.Param("id")
	var task Task
	if result := task.GetTaskDB(id); result.Error != nil {
		c.String(http.StatusNotFound, "A Task with the id '%s' does not exist.", id)
		return
	}
	task.UpdateTaskDB(&params)
	c.JSON(http.StatusOK, task)
}

func handleDeleteTask(c *gin.Context) {

	id := c.Param("id")
	var task Task
	if result := task.GetTaskDB(id); result.Error != nil {
		c.String(http.StatusNotFound, "A Task with the id '%s' does not exist.", id)
		return
	}
	task.DeleteTaskDB()
	c.String(http.StatusOK, "Task removed successfully.")
}

func handleGetTaskStatus(c *gin.Context) {

	id := c.Param("id")
	var task Task
	if result := task.GetTaskDB(id); result.Error != nil {
		c.String(http.StatusNotFound, "A Task with the id '%s' does not exist.", id)
		return
	}
	c.String(http.StatusOK, task.Status)
}

func handlePutTaskStatus(c *gin.Context) {

	var status string
	if err := c.ShouldBindJSON(&status); err != nil {
		c.String(http.StatusBadRequest, "Format not valid")
		return
	}
	if status != "active" && status != "done" {
		c.String(http.StatusBadRequest, "value '%s' is not a legal task status.", status)
		return
	}
	id := c.Param("id")
	var task Task
	if result := task.GetTaskDB(id); result.Error != nil {
		c.String(http.StatusNotFound, "A Task with the id '%s' does not exist.", id)
		return
	}
	task.UpdateTaskDB(&OptParams{Status: status})
	c.String(http.StatusNoContent, "task's status updated successfully.")

}

func handleGetTaskOwnerId(c *gin.Context) {

	id := c.Param("id")
	var task Task
	if result := task.GetTaskDB(id); result.Error != nil {
		c.String(http.StatusNotFound, "A Task with the id '%s' does not exist.", id)
		return
	}
	c.String(http.StatusOK, task.PersonID)
}

func handlePutTaskOwnerId(c *gin.Context) {

	var ownerId string
	if err := c.ShouldBindJSON(&ownerId); err != nil {
		c.String(http.StatusBadRequest, "Format not valid")
		return
	}

	id := c.Param("id")
	var task Task
	if result := task.GetTaskDB(id); result.Error != nil {
		c.String(http.StatusNotFound, "A Task with the id '%s' does not exist.", id)
		return
	}
	task.SetTaskOwnerDB(ownerId)
	c.String(http.StatusNoContent, "task owner updated successfully.")
}
