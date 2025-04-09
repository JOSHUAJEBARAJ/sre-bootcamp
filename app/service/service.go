package service

import (
	"context"

	"github.com/JOSHUAJEBARAJ/sre-bootcamp/models"
	"github.com/JOSHUAJEBARAJ/sre-bootcamp/repository"
)

type StudentService struct {
	repo repository.StudentRepository // this is going to take a interface so we can pass both original and dummy implementation
}

func NewStudentService(r repository.StudentRepository) *StudentService {
	return &StudentService{r}
}

func (s *StudentService) AddStudent(ctx context.Context, input models.StudentInput) (models.Student, error) {
	return s.repo.AddStudent(ctx, input)
}

func (s *StudentService) PingDB(ctx context.Context) error {
	return s.repo.PingDB(ctx)
}

func (s *StudentService) GetStudent(ctx context.Context, id int) (models.Student, error) {
	return s.repo.GetStudent(ctx, id)
}

func (s *StudentService) GetAllStudent(ctx context.Context) ([]models.Student, error) {
	return s.repo.GetAllStudent(ctx)
}

func (s *StudentService) UpdateStudent(ctx context.Context, id int, input models.StudentInput) (models.Student, error) {
	return s.repo.UpdateStudent(ctx, id, input)
}

func (s *StudentService) DeleteStudent(ctx context.Context, id int) error {
	return s.repo.DeleteStudent(ctx, id)
}
