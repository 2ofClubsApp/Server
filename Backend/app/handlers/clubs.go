package handlers

import (
	"../models"
	"../common"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetClubs(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Clubs")
}

func CreateClub(db *gorm.DB, _ http.ResponseWriter, r *http.Request) {
	c := models.NewClub()
	p := common.ExtractPersonInfo(r)
	c.Person = p
	found := common.RecordExists(db, models.ColumnUsername, c.Username, c)
	if !found {
		db.Create(&c)
	}
}

func GetClubsTag(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Clubs Tag")
}

func GetClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

}

func UpdateClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Club")
}
