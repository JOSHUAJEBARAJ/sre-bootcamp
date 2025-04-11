package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/JOSHUAJEBARAJ/sre-bootcamp/handler"
	"github.com/JOSHUAJEBARAJ/sre-bootcamp/models"
	"github.com/JOSHUAJEBARAJ/sre-bootcamp/repository"
	"github.com/JOSHUAJEBARAJ/sre-bootcamp/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	router *gin.Engine
)

func TestMain(m *testing.M) {
	gin.SetMode("test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pgContainer, err := postgres.Run(ctx, "postgres:15.3-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.WithInitScripts("../../db/migrations/000001_create_users_table.up.sql"), // Use your init.sql file
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		panic(err)
	}
	// host, err := pgContainer.Host(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		panic("failed to get container port: " + err.Error())
	}
	dbConfig := models.DatabaseConfig{
		UserName: "testuser",
		Password: "testpass",
		Host:     "localhost",
		Port:     port.Int(),
		DBName:   "testdb",
	}

	db, err := repository.NewDB(ctx, dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	repo, err := repository.NewStudentRepositoryPostgres(db)
	if err != nil {
		log.Fatal(err)
	}
	srv := service.NewStudentService(repo)
	handler := handler.NewStudentHandler(srv)
	router = gin.Default()
	router.GET("/api/v1/students", handler.GetAllStudent)
	router.GET("/api/v1/students/:id", handler.GetStudent)
	router.POST("/api/v1/students", handler.AddStudent)
	code := m.Run()
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			panic("failed to terminate container: " + err.Error())
		}
	}()
	os.Exit(code)

}

func TestAddStudent(t *testing.T) {

	student := models.Student{
		Name:   "Joshua",
		Age:    20,
		Degree: "B.ed",
	}

	// request

	body, _ := json.Marshal(student)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var s models.Student
	respbody, _ := io.ReadAll(w.Body)
	_ = json.Unmarshal(respbody, &s)
	assert.Equal(t, student.Name, s.Name)
}

func TestGetAddStudent(t *testing.T) {

	student := models.Student{
		Name:   "Joshua",
		Age:    20,
		Degree: "B.ed",
	}

	// request

	body, _ := json.Marshal(student)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var s models.Student
	respbody, _ := io.ReadAll(w.Body)
	_ = json.Unmarshal(respbody, &s)

	// get request

	stringID := strconv.Itoa(s.ID)
	getEndpoint := fmt.Sprintf("/api/v1/students/%s", stringID)
	getreq, _ := http.NewRequest(http.MethodGet, getEndpoint, nil)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, getreq)
	getRespBody, _ := io.ReadAll(w.Body)
	var s2 models.Student
	_ = json.Unmarshal(getRespBody, &s2)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, s, s2)

	getEndpoint = "/api/v1/students/10000"
	getreq, _ = http.NewRequest(http.MethodGet, getEndpoint, nil)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, getreq)
	assert.Equal(t, http.StatusNotFound, w.Code)

}
