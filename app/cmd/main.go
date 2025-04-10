package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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
		DBName:   dbName,
	}
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	// TODO figure out way to remove the hardcoded values
	if err := godotenv.Load("/Users/joshua/Projects/sre-bootcamp/.env"); err != nil {
		// DO NOTHING
	}
}

func main() {
	//

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	db, err := repository.NewDB(ctx, populateDBConfig())
	if err != nil {
		log.Fatal(err)
	}
	repo, err := repository.NewStudentRepositoryPostgres(db)
	if err != nil {
		log.Fatal(err)
	}
	srv := service.NewStudentService(repo)
	handler := handler.NewStudentHandler(srv)

	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.GET("/api/v1/students", handler.GetAllStudent)
	router.GET("/api/v1/students/:id", handler.GetStudent)
	router.DELETE("/api/v1/students/:id", handler.DeleteStudent)
	router.PUT("/api/v1/students/:id", handler.PutStudent)
	router.POST("/api/v1/students", handler.AddStudent)
	router.GET("/healthcheck", handler.Healthcheck)

	addr := "localhost:8080"
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Infof("Server running on localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	defer func() {
		db.Close()
	}()

	// handle shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Errorf("Server shutdown failed: %v", err)
	} else {
		log.Info("Server gracefully stopped")
	}

}
