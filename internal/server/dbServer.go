package server

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"sagmi/config"
	"sagmi/internal/models"
)

type DatabaseServer struct {
	DB *gorm.DB
}

func NewDatabaseServer(config *config.Config) *DatabaseServer {
	dbName := config.Viper.GetString("database.name")
	if dbName == "" {
		panic("db name is empty")
	}
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return &DatabaseServer{
		DB: db,
	}
}

func (dbs *DatabaseServer) Start() {
	// Create tables
	err := dbs.DB.AutoMigrate(&models.Endpoint{}, &models.Log{})
	if err != nil {
		log.Fatal(err)
	}
}
