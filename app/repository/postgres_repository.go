package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/JOSHUAJEBARAJ/sre-bootcamp/models"
	log "github.com/sirupsen/logrus"
)

//

type StudentRepositoryPostgres struct {
	// client
	client *sql.DB
}

func NewStudentRepositoryPostgres(ctx context.Context, dbConfig models.DatabaseConfig) (*StudentRepositoryPostgres, error) {
	// todo implemenet this latter

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", dbConfig.Host, dbConfig.Port, dbConfig.UserName, dbConfig.Password, dbConfig.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.WithError(err).Error("Error while opening the connection")
		return &StudentRepositoryPostgres{}, err
	}
	err = db.Ping()
	if err != nil {
		log.WithFields(log.Fields{
			"host": dbConfig.Host,
			"db":   dbConfig.DbName,
		}).WithError(err).Error("Error connecting to the database")
		return &StudentRepositoryPostgres{}, err
	}
	return &StudentRepositoryPostgres{db}, nil
}

func (s *StudentRepositoryPostgres) AddStudent(input models.StudentInput) (models.Student, error) {
	var id int
	err := s.client.QueryRow("INSERT INTO students(name,age,degree) VALUES($1,$2,$3) RETURNING id", input.Name, input.Age, input.Degree).Scan(&id)
	if err != nil {
		log.WithError(err).Error("Error while adding student")
		return models.Student{}, err
	}
	return models.Student{
		Id:     id,
		Name:   input.Name,
		Age:    input.Age,
		Degree: input.Degree,
	}, nil

}

func (s *StudentRepositoryPostgres) GetStudent(id int) (models.Student, error) {
	var student models.Student
	err := s.client.QueryRow("select id, name, age, degree  from students where id=$1", id).Scan(
		&student.Id, &student.Name, &student.Age, &student.Degree)

	if err != nil {
		if err == sql.ErrNoRows {
			// No student found with the given ID
			log.WithField("id", id).Warn("Student not found in DB")
			return models.Student{}, sql.ErrNoRows
		}
		log.WithError(err).Error("Error while getting individual student")
		return models.Student{}, err
	}

	return student, nil
}

func (s *StudentRepositoryPostgres) GetAllStudent() ([]models.Student, error) {
	rows, err := s.client.Query("select id, name, age, degree  from students")
	if err != nil {
		log.WithError(err).Error("Error while getting all student")
		return []models.Student{}, err
	}
	defer rows.Close()
	var students []models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Age, &student.Degree)
		if err != nil {
			log.WithError(err).Error("Error while getting all student")
			return []models.Student{}, err
		}
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		log.WithError(err).Error("Error while getting all student")
		return []models.Student{}, err
	}
	return students, nil
}

func (s *StudentRepositoryPostgres) UpdateStudent(id int, input models.StudentInput) (models.Student, error) {

	var updatedStudent models.Student
	// var id int
	err := s.client.QueryRow(`
	UPDATE students
	set name =$1,age=$2,degree = $3
	where id = $4
	RETURNING id,name,age,degree
	`, input.Name, input.Age, input.Degree, id).Scan(&updatedStudent.Id, &updatedStudent.Name, &updatedStudent.Age, &updatedStudent.Degree)
	if err != nil {
		if err == sql.ErrNoRows {
			log.WithField("id", id).Warn("Student not found in DB")
			return models.Student{}, sql.ErrNoRows
		}
		log.WithError(err).Error("Error while updating")
		return models.Student{}, err
	}
	return updatedStudent, nil

}

func (s *StudentRepositoryPostgres) DeleteStudent(id int) error {

	// var id int
	result, err := s.client.Exec(`
DELETE FROM students where id=$1
	`, id)
	if err != nil {
		log.WithError(err).Error("Error while deleting students")
		return err
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		log.WithError(err).Error("Error while deleting students")
		return err
	}
	if rowAffected == 0 {
		log.WithField("id", id).Warn("Student not found in DB")
		return sql.ErrNoRows
	}
	return nil

}
