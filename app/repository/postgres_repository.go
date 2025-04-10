package repository

import (
	"context"
	"database/sql"

	"github.com/JOSHUAJEBARAJ/sre-bootcamp/models"
	log "github.com/sirupsen/logrus"
)

//

type StudentRepositoryPostgres struct {
	// client
	client *sql.DB
}

func NewStudentRepositoryPostgres(db *DB) (*StudentRepositoryPostgres, error) {
	// todo implemenet this latter

	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s sslmode=disable", dbConfig.Host, dbConfig.Port, dbConfig.UserName, dbConfig.Password, dbConfig.DbName)
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	log.WithError(err).Error("Failed to open connection")
	// 	return nil, err
	// }
	// db.SetMaxOpenConns(MAX_OPEN_CONNECTIONS)
	// db.SetMaxIdleConns(MAX_IDLE_CONNECTIONS)
	// db.SetConnMaxLifetime(time.Minute * MAX_CONNECTION_LIFETIME)
	// err = db.PingContext(ctx)
	// if err != nil {
	// 	log.WithFields(log.Fields{
	// 		"host": dbConfig.Host,
	// 		"db":   dbConfig.DbName,
	// 	}).WithError(err).Error("Failed to ping to the database")
	// 	return nil, err
	// }
	return &StudentRepositoryPostgres{db.client}, nil
}

func (r *StudentRepositoryPostgres) PingDB(ctx context.Context) error {
	// if r.client == nil {
	// 	return errors.New("database client is not initialized")
	// }
	return r.client.PingContext(ctx)
}

func (s *StudentRepositoryPostgres) AddStudent(ctx context.Context, input models.StudentInput) (models.Student, error) {
	var id int
	logger := log.WithFields(log.Fields{
		"name":   input.Name,
		"age":    input.Age,
		"degree": input.Degree,
	})
	logger.Info("Attempting to add new student to database")
	err := s.client.QueryRowContext(ctx, "INSERT INTO students(name,age,degree) VALUES($1,$2,$3) RETURNING id", input.Name, input.Age, input.Degree).Scan(&id)
	if err != nil {
		logger.WithError(err).Error("Failed to add Students")
		return models.Student{}, err
	}
	logger.WithField("student_id", id).Info("Successfully added student to database")
	return models.Student{
		ID:     id,
		Name:   input.Name,
		Age:    input.Age,
		Degree: input.Degree,
	}, nil

}

func (s *StudentRepositoryPostgres) GetStudent(ctx context.Context, id int) (models.Student, error) {
	var student models.Student
	err := s.client.QueryRowContext(ctx, "select id, name, age, degree  from students where id=$1", id).Scan(
		&student.ID, &student.Name, &student.Age, &student.Degree)
	logger := log.WithFields(log.Fields{
		"student_id": id,
	})
	logger.Info("Attempting to get the Student")
	if err != nil {
		if err == sql.ErrNoRows {
			// No student found with the given ID
			logger.Warn("Student not found in DB")
			return models.Student{}, sql.ErrNoRows
		}
		log.WithError(err).Error("Failed to get student")
		return models.Student{}, err
	}

	logger.Info("Successfully found the student")
	return student, nil
}

func (s *StudentRepositoryPostgres) GetAllStudent(ctx context.Context) ([]models.Student, error) {
	rows, err := s.client.QueryContext(ctx, "select id, name, age, degree  from students")
	log.Info("Fetching all students from database")
	if err != nil {
		log.WithError(err).Error("Failed to get all students")
		return []models.Student{}, err
	}
	defer rows.Close()
	var students []models.Student
	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.ID, &student.Name, &student.Age, &student.Degree)
		if err != nil {
			log.WithError(err).Error("Failed to scan student row")
			return []models.Student{}, err
		}
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		log.WithError(err).Error("Error encountered after scanning all rows")
		return []models.Student{}, err
	}
	log.Info("Successfully retrieved all students")
	return students, nil
}

func (s *StudentRepositoryPostgres) UpdateStudent(ctx context.Context, id int, input models.StudentInput) (models.Student, error) {

	var updatedStudent models.Student

	logger := log.WithFields(log.Fields{
		"studnet_id": id,
		"name":       input.Name,
		"age":        input.Age,
		"degree":     input.Degree,
	})

	// var id int
	err := s.client.QueryRowContext(ctx, `
	UPDATE students
	set name =$1,age=$2,degree = $3
	where id = $4
	RETURNING id,name,age,degree
	`, input.Name, input.Age, input.Degree, id).Scan(&updatedStudent.ID, &updatedStudent.Name, &updatedStudent.Age, &updatedStudent.Degree)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn("Student not found in DB")
			return models.Student{}, sql.ErrNoRows
		}
		logger.WithError(err).Error("Failed to update student")
		return models.Student{}, err
	}
	return updatedStudent, nil

}

func (s *StudentRepositoryPostgres) DeleteStudent(ctx context.Context, id int) error {
	logger := log.WithFields(log.Fields{
		"student_id": id,
	})
	// var id int
	logger.Info("Attempting to delete student")
	result, err := s.client.ExecContext(ctx, `
DELETE FROM students where id=$1
	`, id)
	if err != nil {
		logger.WithError(err).Error("Failed to execute delete query")
		return err
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("Failed to get affected rows after deletion")
		return err
	}
	if rowAffected == 0 {
		logger.Warn("Student not found in database")
		return sql.ErrNoRows
	}
	return nil

}
