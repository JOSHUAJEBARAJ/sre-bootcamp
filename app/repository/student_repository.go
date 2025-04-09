package repository

import "github.com/JOSHUAJEBARAJ/sre-bootcamp/models"

type StudentRepository interface {

	// add
	AddStudent(input models.StudentInput) (models.Student, error)
	// get
	GetStudent(id int) (models.Student, error)
	// getAll
	GetAllStudent() ([]models.Student, error)
	//update
	UpdateStudent(id int, input models.StudentInput) (models.Student, error)
	// delete
	DeleteStudent(id int) error
}
