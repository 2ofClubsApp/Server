package app

import (
	"../config"
	"../models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

type App struct {
	db *gorm.DB
}

func (app *App) Initialize(config *config.DBConfig) {
	dbFormat :=
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Host,
			config.Port,
			config.User,
			config.Password,
			config.Name,
		)
	db, err := gorm.Open("postgres", dbFormat)
	if err != nil {
		log.Fatal("Unable to connect to database\n", err)
	}
	defer db.Close()
	app.db = db
	fmt.Println("Connected")
	db.SingularTable(true)
	db.DropTableIfExists(&models.Student{}, &models.Club{}, &models.Event{}, &models.Chat{}, &models.Log{})
	db.CreateTable(&models.Student{}, &models.Club{}, &models.Event{}, &models.Chat{}, &models.Log{})
}
