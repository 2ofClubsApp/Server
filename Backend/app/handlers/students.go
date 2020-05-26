package handlers

import (
	"../models"
	"../common"
	"../models/status"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

func CreateStudent(db *gorm.DB, _ http.ResponseWriter, r *http.Request) {
	// Check if content type is application/json?
	s := models.NewStudent()
	p := common.ExtractPersonInfo(r)
	s.Person = p
	found := common.RecordExists(db, models.ColumnUsername, s.Username, s)
	if !found {
		db.Create(&s)
	}
}

func GetStudent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars[models.ColumnUsername]
	s := models.NewStudent()
	ss := status.New()
	// Defaults will be overridden when obtaining data and being inserted into struct except for null
	found := common.RecordExists(db, models.ColumnUsername, username, s)
	if !found {
		ss.Message = status.FAILURE
	} else {
		ss.Data = s
	}
	common.WriteData(common.ParseJSON(ss), http.StatusOK, w)
}

func UpdateStudent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Student")
}
