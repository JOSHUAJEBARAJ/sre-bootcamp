package repository

import (
	"context"

	"github.com/JOSHUAJEBARAJ/sre-bootcamp/models"
)

type StudentRepository interface {
	// add
	AddStudent(ctx context.Context, input models.StudentInput) (models.Student, error)
	// get
	GetStudent(ctx context.Context, id int) (models.Student, error)
	// getAll
	GetAllStudent(ctx context.Context) ([]models.Student, error)
	//update
	UpdateStudent(ctx context.Context, id int, input models.StudentInput) (models.Student, error)
	// delete
	DeleteStudent(ctx context.Context, id int) error
	// TODO : Move this to the seperate interface
	PingDB(ctx context.Context) error
}
