package service

import (
	"github.com/JOSHUAJEBARAJ/sre-bootcamp/models"
	"github.com/JOSHUAJEBARAJ/sre-bootcamp/repository"
)

type StudentService struct {
	repo repository.StudentRepository // this is going to take a interface so we can pass both original and dummy implementation
}

func NewStudentService(r repository.StudentRepository) *StudentService {
	return &StudentService{r}
}

func (s *StudentService) AddStudent(input models.StudentInput) (models.Student, error) {
	return s.repo.AddStudent(input)
}

func (s *StudentService) GetStudent(id int) (models.Student, error) {
	return s.repo.GetStudent(id)
}

func (s *StudentService) GetAllStudent() ([]models.Student, error) {
	return s.repo.GetAllStudent()
}

func (s *StudentService) UpdateStudent(id int, input models.StudentInput) (models.Student, error) {
	return s.repo.UpdateStudent(id, input)
}

func (s *StudentService) DeleteStudent(id int) error {
	return s.repo.DeleteStudent(id)
}
