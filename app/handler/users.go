package handler

import (
	"../model"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

/*
Create a User given a JSON payload. (See models.User for payload information).
*/
func CreateUser(db *gorm.DB, w http.ResponseWriter, c *model.Credentials, u *model.User) {
	u.Credentials = c
	status := model.NewStatus()
	status.Message = SignupSuccess
	db.Create(&u)
	WriteData(GetJSON(status), http.StatusOK, w)
}

/*
Validating the user request to ensure that they can only access/modify their own data.
If valid, (sub, true) is returned, otherwise (sub, false) where sub represents the username accessing the resource.
*/
func IsValidRequest(username string, r *http.Request) bool {
	claims := GetTokenClaims(r)
	sub := fmt.Sprintf("%v", claims["sub"])
	fmt.Println(username)
	fmt.Println(sub)
	return sub == username
}

// FIX: Extract user Club info from "Manages" (List them as club names rather than the entire club? with an isOwner) in JSON
func GetUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var data string
	vars := mux.Vars(r)
	username := strings.ToLower(vars[model.UsernameColumn])
	status := model.NewStatus()
	u := model.NewUser()

	if IsValidRequest(username, r) {
		// Defaults will be overridden when obtaining data and being inserted into struct except for null
		found := SingleRecordExists(db, model.UserTable, model.UsernameColumn, username, u)
		db.Table(model.UserTable).Preload(model.ManagesColumn).Find(u)
		if !found {
			status.Message = model.UserNotFound
			status.Code = model.FailureCode
		} else {
			status.Data = u
			status.Message = model.UserFound
		}
		statusCode = http.StatusOK
	} else {
		status.Message = http.StatusText(http.StatusForbidden)
		statusCode = http.StatusForbidden
		status.Code = -1
	}
	data = GetJSON(status)
	WriteData(data, statusCode, w)
}

func UpdateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update User")
}
