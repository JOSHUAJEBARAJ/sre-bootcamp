package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/JOSHUAJEBARAJ/sre-bootcamp/handler"
	"github.com/JOSHUAJEBARAJ/sre-bootcamp/models"
	"github.com/JOSHUAJEBARAJ/sre-bootcamp/repository"
	"github.com/JOSHUAJEBARAJ/sre-bootcamp/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func populateDBConfig() models.DatabaseConfig {
	username := os.Getenv("DB_USERNAME")
	if username == "" {
		log.Fatal("DB_USERNAME SHOULD BE SET")
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		log.Fatal("DB_PASSWORD SHOULD BE SET")
	}
	host := os.Getenv("DB_HOST")
	if host == "" {
		log.Fatal("DB_HOST SHOULD BE SET")
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		log.Fatal("DB_PORT SHOULD BE SET")
	}
	intPort, err := strconv.ParseInt(port, 0, 0)
	if err != nil {
		log.Fatal("Error while converting the port number", err)
	}
	dbName := os.Getenv("DB_DBNAME")
	if dbName == "" {
		log.Fatal("DB_DBNAME SHOULD BE SET")
	}
	return models.DatabaseConfig{
		UserName: username,
		Password: password,
		Host:     host,
		Port:     int(intPort),
		DbName:   dbName,
	}
}
func main() {
	//

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}
	db := populateDBConfig()
	repo, err := repository.NewStudentRepositoryPostgres(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	srv := service.NewStudentService(repo)
	handler := handler.NewStudentHandler(srv)

	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.GET("/students", handler.GetAllStudent)
	router.GET("/students/:id", handler.GetStudent)
	router.DELETE("/students/:id", handler.DeleteStudent)
	router.PUT("/students/:id", handler.PutStudent)
	router.POST("/students", handler.AddStudent)
	port := os.Getenv("SERVER_PORT")

	if port == "" {
		port = "8080"
	}
	log.Infof("Server running on localhost:%s", port)
	err = router.Run("localhost:" + port)
	if err != nil {
		log.Fatal("Error while Creating the server", err)
	}

}
