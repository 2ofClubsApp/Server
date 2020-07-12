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

func CreateClub(db *gorm.DB, w http.ResponseWriter, p *model.User, c *model.Club, tableName string) {
	fmt.Println("Creating a Club")
}

func GetClubsTag(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Clubs Tag")
}

func GetClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting club")
}

func UpdateClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Club")
}
