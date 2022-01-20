package handlers

import (
	. "backend/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetPeople(c *gin.Context) {
	var people []Person
	var person Person
	person.GetPeopleDB(&people)
	c.JSON(http.StatusOK, people)
}

func HandleGetPerson(c *gin.Context) {
	if c.Request.Method == "GET" {
		id := c.Param("id")
		var person Person
		if result := person.GetPersonDB(id); result.Error != nil {
			c.String(http.StatusNotFound, ValidateDbErrorMsg(result.Error, "person", &OptParams{ID: id}))

			return
		}
		c.JSON(http.StatusOK, person)
	}
}

func HandlePatchPerson(c *gin.Context) {

	create := struct {
		Name                        *string `json:"name" binding:"omitempty,notemptystring"`
		Email                       *string `json:"email" binding:"omitempty,notemptystring"`
		FavoriteProgrammingLanguage *string `json:"favoriteProgrammingLanguage" binding:"omitempty,notemptystring"`
	}{}
	if err := c.ShouldBindJSON(&create); err != nil {
		c.String(http.StatusBadRequest, ValidateTagErrorMsg(err))
		return
	}

	id := c.Param("id")

	var person Person
	if result := person.GetPersonDB(id); result.Error != nil {
		c.String(http.StatusNotFound, ValidateDbErrorMsg(result.Error, "person", &OptParams{ID: id}))
		return
	}
	result := person.UpdatePersonDB(&OptParams{
		Name:                        create.Name,
		Email:                       create.Email,
		FavoriteProgrammingLanguage: create.FavoriteProgrammingLanguage,
	})

	if result.Error != nil {
		c.String(http.StatusBadRequest, ValidateDbErrorMsg(result.Error, "person", &OptParams{Email: &person.Email}))

		return
	}
	c.JSON(http.StatusOK, person)
}

func HandleDeletePerson(c *gin.Context) {
	var person Person
	id := c.Param("id")
	if result := person.GetPersonDB(id); result.Error != nil {
		c.String(http.StatusNotFound, ValidateDbErrorMsg(result.Error, "person", &OptParams{ID: id}))
		return
	}
	person.DeletePersonDB(id)
	c.String(http.StatusOK, "Person removed successfully.")
}

func HandlePostPerson(c *gin.Context) {
	var person Person
	create := struct {
		Name                        string `json:"name" binding:"required"`
		Email                       string `json:"email" binding:"required,email"`
		FavoriteProgrammingLanguage string `json:"favoriteProgrammingLanguage" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&create); err != nil {

		c.String(http.StatusBadRequest, ValidateTagErrorMsg(err))
		return
	}
	result := person.CreatePersonDB(&OptParams{Name: &create.Name, Email: &create.Email, FavoriteProgrammingLanguage: &create.FavoriteProgrammingLanguage})
	if result.Error != nil {
		c.String(http.StatusBadRequest, ValidateDbErrorMsg(result.Error, "person", &OptParams{Email: &person.Email}))

		return

	}

	c.Header("Location", fmt.Sprintf("/api/people/%s", person.ID))
	c.Header("x-Created-Id", fmt.Sprint(person.ID))
	c.String(http.StatusCreated, "Person created successfully")
}
