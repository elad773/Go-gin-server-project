package model
import "time"

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
	ID                          string
	Name                        *string `json:"name"`
	Email                       *string `json:"email"`
	FavoriteProgrammingLanguage *string `json:"favoriteProgrammingLanguage"`
	Title                       *string `json:"title"`
	Details                     *string `json:"details"`
	DueDate                     *string `json:"dueDate"`
	Status                      string  `json:"status"`
}