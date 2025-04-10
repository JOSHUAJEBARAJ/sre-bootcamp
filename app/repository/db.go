package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JOSHUAJEBARAJ/sre-bootcamp/models"
	log "github.com/sirupsen/logrus"
)

const (
	MAX_OPEN_CONNECTIONS    = 25
	MAX_IDLE_CONNECTIONS    = 25
	MAX_CONNECTION_LIFETIME = 5
)

type DB struct {
	client *sql.DB
}

func NewDB(ctx context.Context, config models.DatabaseConfig) (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.UserName, config.Password, config.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.WithError(err).Error("Failed to open connection")
		return nil, err
	}
	db.SetMaxOpenConns(MAX_OPEN_CONNECTIONS)
	db.SetMaxIdleConns(MAX_IDLE_CONNECTIONS)
	db.SetConnMaxLifetime(time.Minute * MAX_CONNECTION_LIFETIME)
	if err := db.PingContext(ctx); err != nil {
		log.WithFields(log.Fields{
			"host": config.Host,
			"db":   config.DBName,
		}).WithError(err).Error("Failed to ping to the database")
		return nil, err
	}
	return &DB{client: db}, nil
}

func (d *DB) Close() error {
	return d.client.Close()
}

func (d *DB) Ping(ctx context.Context) error {
	if err := d.client.PingContext(ctx); err != nil {
		log.WithError(err).Error("Database ping failed")
		return err
	}
	return nil
}
