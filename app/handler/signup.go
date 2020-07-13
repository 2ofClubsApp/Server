package handler

import (
	"../model"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func isCredAvailable(db *gorm.DB, r *http.Request, tableName, column string) bool {
	var placeholder interface{}
	vars := mux.Vars(r)
	cred := vars[column]
	return !RecordExists(db, tableName, column, cred, placeholder)
}

/*
Returning the availability of a username or email.
On success nothing is return otherwise, a status error is returned
 */
func returnRequest(fieldAvailable bool, w http.ResponseWriter, response string) int {
	if fieldAvailable {
		s := model.NewStatus()
		s.Message = response
		data := GetJSON(s)
		return WriteData(data, http.StatusOK, w)
	}
	return -1
}

/*
Querying username against database to find the availability of username
 */
func QueryUsername(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	unameAvailable := isCredAvailable(db, r, model.UserTable, model.UsernameColumn)
	returnRequest(unameAvailable, w, model.UsernameExists)
}

func QueryEmail(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	emailAvailable := isCredAvailable(db, r, model.UserTable, model.EmailColumn)
	returnRequest(emailAvailable, w, model.EmailExists)
}
