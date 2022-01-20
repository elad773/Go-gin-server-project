package handlers

import (
	. "backend/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HandlePostTask(c *gin.Context) {
	id := c.Param("id")
	var person Person
	var task Task
	create := struct {
		Title   string `json:"title" binding:"required"`
		Details string `json:"details" binding:"required"`
		DueDate string `json:"dueDate" binding:"required"`
		Status  string `json:"status"`
	}{}
	if err := c.ShouldBindJSON(&create); err != nil {
		c.String(http.StatusBadRequest, ValidateTagErrorMsg(err))
		return
	}
	if _, dateErr := time.Parse("2006-01-02", create.DueDate); dateErr != nil {
		if _, timeErr := time.Parse(time.RFC3339, create.DueDate); timeErr != nil {
			c.String(http.StatusBadRequest, "dueDate '%s' is not in RFC3339 date format.", create.DueDate)

			return
		}
	}

	if result := person.GetPersonDB(id); result.Error != nil {
		c.String(http.StatusNotFound, ValidateDbErrorMsg(result.Error, "person", &OptParams{ID: id}))
		return
	}
	err := task.CreateTaskDB(&person, &OptParams{Title: &create.Title, Details: &create.Details, DueDate: &create.DueDate, Status: create.Status})
	if err != nil {
		c.String(http.StatusBadRequest, ValidateDbErrorMsg(err, "task", &OptParams{}))
		return
	}
	c.Header("Location", fmt.Sprintf("/api/tasks/%s", task.ID))
	c.Header("x-Created-Id", fmt.Sprint(task.ID))
	c.String(http.StatusCreated, "Task created successfuly.")

}

func HandleGetTasks(c *gin.Context) {

	id := c.Param("id")
	status := c.Query("status")
	var person Person
	if result := person.GetPersonDB(id); result.Error != nil {
		c.String(http.StatusNotFound, ValidateDbErrorMsg(result.Error, "person", &OptParams{ID: id}))
		return
	}
	person.GetTasksDB(status)
	c.JSON(http.StatusOK, person.Tasks)

}

func HandleGetTask(c *gin.Context) {
	id := c.Param("id")
	var task Task
	if result := task.GetTaskDB(id); result.Error != nil {
		c.String(http.StatusNotFound, ValidateDbErrorMsg(result.Error, "task", &OptParams{ID: id}))
		return
	}
	c.JSON(http.StatusOK, task)
}

func HandlePatchTask(c *gin.Context) {
	create := struct {
		Title   *string `json:"title" binding:"omitempty,notemptystring"`
		Details *string `json:"details"`
		DueDate *string `json:"dueDate"`
		Status  string  `json:"status"`
	}{}
	if err := c.ShouldBindJSON(&create); err != nil {
		c.String(http.StatusBadRequest, ValidateTagErrorMsg(err))
		return
	}
	if create.Details != nil {
		if _, dateErr := time.Parse("2006-01-02", *create.DueDate); dateErr != nil {
			if _, timeErr := time.Parse(time.RFC3339, *create.DueDate); timeErr != nil {
				c.String(http.StatusBadRequest, "dueDate '%s' is not in RFC3339 date format.", create.DueDate)

				return
			}
		}
	}
	id := c.Param("id")
	var task Task
	if result := task.GetTaskDB(id); result.Error != nil {
		c.String(http.StatusNotFound, ValidateDbErrorMsg(result.Error, "task", &OptParams{ID: id}))
		return
	}
	task.UpdateTaskDB(&OptParams{Title: create.Title, Details: create.Details, DueDate: create.DueDate, Status: create.Status})
	c.JSON(http.StatusOK, task)
}

func HandleDeleteTask(c *gin.Context) {

	id := c.Param("id")
	var task Task
	if result := task.GetTaskDB(id); result.Error != nil {
		c.String(http.StatusNotFound, ValidateDbErrorMsg(result.Error, "task", &OptParams{ID: id}))
		return
	}
	task.DeleteTaskDB()
	c.String(http.StatusOK, "Task removed successfully.")
}

func HandleGetTaskStatus(c *gin.Context) {

	id := c.Param("id")
	var task Task
	if result := task.GetTaskDB(id); result.Error != nil {
		c.String(http.StatusNotFound, ValidateDbErrorMsg(result.Error, "task", &OptParams{ID: id}))
		return
	}
	c.String(http.StatusOK, task.Status)
}

func HandlePutTaskStatus(c *gin.Context) {

	var status string
	if err := c.ShouldBindJSON(&status); err != nil {
		c.String(http.StatusBadRequest, "Format not valid. Check this error:\n%s", err.Error())
		return
	}
	if status != "active" && status != "done" {
		c.String(http.StatusBadRequest, "value '%s' is not a legal task status.", status)
		return
	}
	id := c.Param("id")
	var task Task
	if result := task.GetTaskDB(id); result.Error != nil {
		c.String(http.StatusNotFound, ValidateDbErrorMsg(result.Error, "task", &OptParams{ID: id}))
		return
	}
	task.UpdateTaskDB(&OptParams{Status: status})
	c.String(http.StatusNoContent, "task's status updated successfully.")

}

func HandleGetTaskOwnerId(c *gin.Context) {

	id := c.Param("id")
	var task Task
	if result := task.GetTaskDB(id); result.Error != nil {
		c.String(http.StatusNotFound, ValidateDbErrorMsg(result.Error, "task", &OptParams{ID: id}))
		return
	}
	c.String(http.StatusOK, task.PersonID)
}

func HandlePutTaskOwnerId(c *gin.Context) {
	var ownerId string
	if err := c.ShouldBindJSON(&ownerId); err != nil {
		c.String(http.StatusBadRequest, "Format not valid. Check this error:\n%s", err.Error())
		return
	}
	id := c.Param("id")
	var task Task
	if result := task.GetTaskDB(id); result.Error != nil {
		c.String(http.StatusNotFound, ValidateDbErrorMsg(result.Error, "task", &OptParams{ID: id}))
		return
	}
	if err := task.SetTaskOwnerDB(ownerId); err != nil {
		c.String(http.StatusNotFound, ValidateDbErrorMsg(err, "person", &OptParams{ID: ownerId}))
		return
	}
	c.String(http.StatusNoContent, "task owner updated successfully.")
}
