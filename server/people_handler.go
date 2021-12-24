package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleGetPeople(c *gin.Context) {
	var people []Person
	var person Person 
	person.GetPeopleDB(&people)
	c.JSON(http.StatusOK, people)
}

func handleGetPerson(c *gin.Context) {

	id := c.Param("id")
	var person Person
	if result := person.GetPersonDB(id); result.Error != nil {
		c.String(http.StatusNotFound, "A person with the id '%s'does not exist. ", id)
		return
	}
	c.JSON(http.StatusOK, person)

}

func handlePatchPerson(c *gin.Context) {

	create := struct {
		Name                        *string `json:"name"`
		Email                       *string `json:"email"`
		FavoriteProgrammingLanguage *string `json:"favoriteProgrammingLanguage"`
	}{}
	if err := c.ShouldBindJSON(&create); err != nil {

		return
	}

	id := c.Param("id")
	var person Person
	if result := person.GetPersonDB(id); result.Error != nil {
		c.String(http.StatusNotFound, "A person with the id '%s'does not exist. ", id)
		return
	}
	result := person.UpdatePersonDB( &OptParams{
		Name:                        create.Name,
		Email:                       create.Email,
		FavoriteProgrammingLanguage: create.FavoriteProgrammingLanguage,
	})
	if result.Error != nil {
		if result.Error.Error() == "UNIQUE constraint failed: people.email" {
			c.String(http.StatusBadRequest, "A person with email '%s' already exists.", person.Email)
			return
		}
	}
	c.JSON(http.StatusOK, person)
}

func handleDeletePerson(c *gin.Context) {
    var person Person
	id := c.Param("id")
	if result := person.GetPersonDB(id); result.Error != nil {
		c.String(http.StatusNotFound, "A person with the id '%s' does not exist.", id)
		return
	}
	person.DeletePersonDB(id)
	c.String(http.StatusOK, "Person removed successfully.")
}

func handlePostPerson(c *gin.Context) {
	var person Person
	create := struct {
		Name                        string `json:"name" binding:"required"`
		Email                       string `json:"email" binding:"required"`
		FavoriteProgrammingLanguage string `json:"favoriteProgrammingLanguage" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&create); err != nil {
		c.String(http.StatusBadRequest, "%s", "Required data fields are missing, data makes no sense, or data contains illegal values.\n"+err.Error())
		return
	}
	result :=person.CreatePersonDB(&OptParams{Name: &create.Name, Email: &create.Email, FavoriteProgrammingLanguage: &create.FavoriteProgrammingLanguage})
	if result.Error != nil {
		if result.Error.Error() == "UNIQUE constraint failed: people.email" {
			c.String(http.StatusBadRequest, "A person with email '%s' already exists.", person.Email)
		} else {
			c.String(http.StatusBadRequest, "Required data fields are missing, data makes no sense, or data contains illegal values.\n"+result.Error.Error())
		}
		return
	}

	c.Header("Location", fmt.Sprintf("/api/people/%s", person.ID))
	c.Header("x-Created-Id", fmt.Sprint(person.ID))
	c.String(http.StatusCreated, "Person created successfully")
}
