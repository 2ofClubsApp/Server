package handler

import (
	"../model"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetClubs(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Clubs")
}

func CreateClub(db *gorm.DB, _ http.ResponseWriter, r *http.Request) {
	c := model.NewClub()
	p := ExtractPersonInfo(r)
	c.Person = p
	found := RecordExists(db, model.ClubTable, model.UsernameColumn, c.Username, c)
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