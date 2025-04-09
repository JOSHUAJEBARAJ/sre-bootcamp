package models

// model for storing the Student

type Student struct {
	Id     int
	Name   string
	Age    int
	Degree string
}

// student Input

type StudentInput struct {
	Name   string
	Age    int
	Degree string
}

// database config

type DatabaseConfig struct {
	// We will fill this later
	Host     string
	Port     int
	UserName string
	Password string
	DbName   string
}
