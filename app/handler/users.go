package handler

import (
	"../model"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

/*
Create a User given a JSON payload. (See models.User for payload information).
*/
func CreateUser(db *gorm.DB, w http.ResponseWriter, c *model.Credentials, u *model.User) {
	u.Credentials = c
	status := model.NewStatus()
	db.Create(&u)
	WriteData(GetJSON(status), http.StatusOK, w)
}

func GetUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	data := ""
	vars := mux.Vars(r)
	username := vars[model.UsernameColumn]
	if ValidateUserReq(username, r) {
		s := model.NewUser()
		ss := model.NewStatus()
		// Defaults will be overridden when obtaining data and being inserted into struct except for null
		found := RecordExists(db, "user", model.UsernameColumn, username, s)
		if !found {
			ss.Message = model.UserNotFound
		} else {
			ss.Data = s
		}
		data = GetJSON(ss)
	} else {
		data = http.StatusText(http.StatusForbidden)
		status = http.StatusForbidden
	}
	WriteData(data, status, w)

}

func UpdateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update User")
}
