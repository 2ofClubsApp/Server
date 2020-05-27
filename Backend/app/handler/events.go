package handler

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetEvents(db *gorm.DB, w http.ResponseWriter, r *http.Request){
	fmt.Println("Get Events")
}

func GetEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request){
	fmt.Println("Get Event")
}

func CreateEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request){
	fmt.Println("Create Event")
}


func UpdateEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request){
	fmt.Println("Update Event")
}

func DeleteEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request){
	fmt.Println("Delete Event")
}
