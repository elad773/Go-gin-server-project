package server

import (
	"time"

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

//struct for person data. Use for data communication between DB and server.
type Person struct {
	CreatedAt                   time.Time  `json:"-"`
	UpdatedAt                   time.Time  `json:"-"`
	DeletedAt                   *time.Time `json:"-" gorm:"index"`
	Name                        string     `json:"name"`
	Email                       string     `json:"email" gorm:"unique"`
	FavoriteProgrammingLanguage string     `json:"favoriteProgrammingLanguage"`
	ActiveTaskCount             uint       `json:"activeTaskCount"`
	ID                          string     `json:"id" gorm:"primarykey"`
	Tasks                       []Task     `json:"-"`
}

//struct for task data. Use for data communication between DB and server.
type Task struct {
	ID        string     `json:"-" gorm:"primarykey"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
	Title     string     `json:"title"`
	Details   string     `json:"details"`
	DueDate   string     `json:"dueDate"`
	Status    string     `json:"status"`
	PersonID  string     `json:"ownerId"`
}

//struct that store optional data of person and task
// It purpose is to support optional args in functions (as go not support that).
type OptParams struct {
	Name                        *string `json:"name"`
	Email                       *string `json:"email"`
	FavoriteProgrammingLanguage *string `json:"favoriteProgrammingLanguage"`
	Title                       *string `json:"title"`
	Details                     *string `json:"details"`
	DueDate                     *string `json:"dueDate"`
	Status                      string  `json:"status"`
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

func (person *Person) CreatePersonDB(optParams *OptParams) (tx *gorm.DB) {
	person.Name = *optParams.Name
	person.Email = *optParams.Email
	person.FavoriteProgrammingLanguage = *optParams.FavoriteProgrammingLanguage
	person.ActiveTaskCount = 0
	person.Tasks = []Task{}
	return db.Create(&person)

}
func (person *Person) GetPersonDB(id string) (tx *gorm.DB) {
	//gorm.DB.First - query for the first row that satisfy the condition
	return db.Where("ID = ?", id).First(&person)
}
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
func (person *Person) DeletePersonDB(id string) {
	db.Where("ID = ?", id).Delete(&person)
}
func (task *Task) CreateTaskDB(person *Person, params *OptParams) (err error) {
	task.PersonID = person.ID
	task.Title = *params.Title
	task.Details = *params.Details
	task.DueDate = *params.DueDate
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

func (person *Person) GetTasksDB(status string) {
	if status == "active" || status == "done" {
		db.Preload("Tasks", "Status = ?", status).First(&person)
	} else {
		db.Preload("Tasks").First(&person)
	}
}
func (task *Task) GetTaskDB(id string) (tx *gorm.DB) {
	return db.Where("ID = ?", id).First(&task)
}
func (task *Task) UpdateTaskDB(params *OptParams) {

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

	db.Save(&task)
}
func (task *Task) DeleteTaskDB() {
	var person Person
	person.GetPersonDB(task.PersonID)
	person.ActiveTaskCount--
	person.UpdatePersonDB(&OptParams{})
	db.Delete(&task)
}
func (task *Task) SetTaskOwnerDB(ownerId string) {
	if task.Status == "active" {
		var person Person
		var person1 Person
		db.Where("ID = ?", task.PersonID).First(&person)
		person.ActiveTaskCount = person.ActiveTaskCount - 1
		db.Save(&person)
		db.Where("ID = ?", ownerId).First(&person1)
		person1.ActiveTaskCount = person1.ActiveTaskCount + 1
		db.Save(&person1)

	}
	task.PersonID = ownerId
	db.Save(&task)
}
