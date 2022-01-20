package model

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func OpenDB() {
	create_db, err := gorm.Open(sqlite.Open("sqlite1.db"), &gorm.Config{})
	if err != nil {

	}
	db = create_db
	db.AutoMigrate(&Person{}, &Task{})
}

//create unique id for the new created person row in DB.
func (p *Person) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewString()
	return

}

//create unique id for the new created task row in DB.
func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.NewString()
	return

}

//create new column person in DB of Person data.
func (person *Person) CreatePersonDB(optParams *OptParams) (tx *gorm.DB) {
	person.Name = *optParams.Name
	person.Email = *optParams.Email
	person.FavoriteProgrammingLanguage = *optParams.FavoriteProgrammingLanguage
	person.ActiveTaskCount = 0
	person.Tasks = []Task{}
	return db.Create(&person)

}

//retrive Person data from DB, into person
func (person *Person) GetPersonDB(id string) (tx *gorm.DB) {
	//gorm.DB.First - query for the first row that satisfy the condition
	return db.Where("ID = ?", id).First(&person)
}

//retrive Person data of all the Persons from DB, into people
func (person *Person) GetPeopleDB(people *[]Person) (tx *gorm.DB) {
	return db.Find(&people)
}

//updte person fields that were included in the json.
// if field were not included, it not change in DB.
func (person *Person) UpdatePersonDB(optParams *OptParams) (tx *gorm.DB) {
	if optParams.Name != nil {
		person.Name = *optParams.Name
	}
	if optParams.Email != nil {
		person.Email = *optParams.Email
	}
	if optParams.FavoriteProgrammingLanguage != nil {
		person.FavoriteProgrammingLanguage = *optParams.FavoriteProgrammingLanguage
	}
	return db.Save(&person)
}

//delete Person data in DB of person
func (person *Person) DeletePersonDB(id string) {
	db.Preload("Tasks").First(&person)
	for _, task := range person.Tasks {
		task.DeleteTaskDB()
	}
	//db.Where("ID = ?", id).Delete(&person)
	db.Delete(&person)
}

//create new column Task in DB of task data.
func (task *Task) CreateTaskDB(person *Person, params *OptParams) (err error) {
	task.PersonID = person.ID
	task.Title = *params.Title
	task.Details = *params.Details
	task.DueDate = *params.DueDate
	if params.Status == "" {
		params.Status = "active"
	}
	if params.Status != "" && params.Status != "active" && params.Status != "done" {
		return errors.New("Status data is Invalid.")
	}
	task.Status = params.Status
	result := db.Create(&task)
	if result.Error != nil {
		return result.Error
	}
	if task.Status == "active" {
		person.ActiveTaskCount++
		db.Save(&person)
	}
	return
}

//retrive Tasks data of person from DB, into Tasks field in person
func (person *Person) GetTasksDB(status string) {
	if status == "active" || status == "done" {
		db.Preload("Tasks", "Status = ?", status).First(&person)
	} else {
		db.Preload("Tasks").First(&person)
	}
}

//retrive Task data from DB, into task
func (task *Task) GetTaskDB(id string) (tx *gorm.DB) {
	return db.Where("ID = ?", id).First(&task)
}

//updte task fields that were included in the json or that included for update.
// if field were not included, it not change in DB.
func (task *Task) UpdateTaskDB(params *OptParams) (err error) {

	if params.Title != nil {
		task.Title = *params.Title
	}
	if params.Details != nil {
		task.Details = *params.Details
	}
	if params.DueDate != nil {
		task.DueDate = *params.DueDate
	}
	if params.Status == "active" || params.Status == "done" {
		var person Person
		person.GetPersonDB(task.PersonID)
		if params.Status == "done" && task.Status == "active" {
			person.ActiveTaskCount--
		}
		if params.Status == "active" && task.Status == "done" {
			person.ActiveTaskCount++
		}
		db.Save(&person)
		task.Status = params.Status
	}

	return db.Save(&task).Error
}

//delete Task data in DB of task
func (task *Task) DeleteTaskDB() {
	var person Person
	person.GetPersonDB(task.PersonID)
	person.ActiveTaskCount--
	person.UpdatePersonDB(&OptParams{})
	db.Delete(&task)
}

// update the OwenerId field data in DB of task
func (task *Task) SetTaskOwnerDB(ownerId string) (err error) {
	var person1 Person
	var person Person
	fmt.Println(task.Status)
	fmt.Println(task.DueDate)
	if task.Status == "active" {

		result := person1.GetPersonDB(ownerId)
		if result.Error != nil {
			return result.Error
		}
		person.GetPersonDB(task.PersonID)
		person.ActiveTaskCount = person.ActiveTaskCount - 1
		db.Save(&person)
		person1.ActiveTaskCount = person1.ActiveTaskCount + 1
		db.Save(&person1)

	}
	task.PersonID = ownerId
	db.Save(&task)
	return nil
}
