package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ElayadeIsmail/go-pingram/config"
	"github.com/ElayadeIsmail/go-pingram/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	p := config.Getenv("DB_PORT")
	port, _ := strconv.ParseInt(p, 10, 32)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Getenv("DB_HOST"), config.Getenv("DB_USER"), config.Getenv("DB_PASSWORD"), config.Getenv("DB_NAME"), port)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could Not Connect To Database")
	}
	fmt.Println("Connected Successfully to database")
	DB.AutoMigrate(&models.User{}, &models.Post{})
	fmt.Println("Database Migrated")
}

func Close() {
	pgDB, err := DB.DB()
	if err != nil {
		log.Fatalln(err.Error())
	}
	pgDB.Close()
}
