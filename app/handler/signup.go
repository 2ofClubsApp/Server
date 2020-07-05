package handler

import (
	"../model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

func isUserAvailable(db *gorm.DB, r *http.Request, tableName, column string) bool {
	var placeholder interface{}
	switch tableName {
	case model.ClubTable:
		placeholder = model.NewClub()
		break
	case model.StudentTable:
		placeholder = model.NewStudent()
	}
	vars := mux.Vars(r)
	email := vars[column]
	return !RecordExists(db, tableName, column, email, placeholder)
}

func returnRequest(fparam bool, sparam bool, w http.ResponseWriter, response string) int {
	if !(fparam && sparam) {
		s := model.NewStatus()
		s.Message = response
		data := ParseJSON(s)
		return WriteData(data, http.StatusOK, w)
	}
	return -1
}

func QueryUsername(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	userClub := isUserAvailable(db, r, model.ClubTable, model.UsernameColumn)
	userStudent := isUserAvailable(db, r, model.StudentTable, model.UsernameColumn)
	returnRequest(userClub, userStudent, w, model.UsernameFound)
}

func QueryEmail(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	emailClub := isUserAvailable(db, r, model.ClubTable, model.EmailColumn)
	emailStudent := isUserAvailable(db, r, model.StudentTable, model.EmailColumn)
	returnRequest(emailClub, emailStudent, w, model.EmailFound)
}
