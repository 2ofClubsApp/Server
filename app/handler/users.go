package handler

import (
	"../model"
	"encoding/json"
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
	db.Create(u)
	WriteData(GetJSON(status), http.StatusOK, w)
}

/*
Validating the user request to ensure that they can only access/modify their own data.
True if the requested user has the same username identifier as the token username
*/
func IsValidRequest(username string, r *http.Request) bool {
	claims := GetTokenClaims(r)
	sub := fmt.Sprintf("%v", claims["sub"])
	//fmt.Println(username)
	//fmt.Println(sub)
	return sub == username
}

// FIX: Extract user Club info from "Manages" (List them as club names with an isOwner) in JSON
func GetUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var httpStatus int
	var data string
	vars := mux.Vars(r)
	username := strings.ToLower(vars["username"])
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
		httpStatus = http.StatusOK
	} else {
		status.Message = http.StatusText(http.StatusForbidden)
		httpStatus = http.StatusForbidden
		status.Code = -1
	}
	data = GetJSON(status)
	WriteData(data, httpStatus, w)
}

/*
Updating the users choice of tags and attended events. Only valid tags will be extracted and added if it's not already.
If an invalid format is provided where there aren't any valid tags to be extracted, the users tag preferences will be reset
*/
func UpdateUserTags(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var httpStatus int
	status := model.NewStatus()
	vars := mux.Vars(r)
	username := strings.ToLower(vars["username"])
	if IsValidRequest(username, r) {
		user := model.NewUser()
		// User is guaranteed to have an account (Verified JWT and request is verified)
		SingleRecordExists(db, model.UserTable, model.UsernameColumn, username, user)
		var chooses []*model.Tag
		for _, name := range getTagInfo(r) {
			tag := model.NewTag()
			if SingleRecordExists(db, model.TagTable, model.NameColumn, name, tag) {
				tag.Name = name
				chooses = append(chooses, tag)
			}
		}
		db.Model(user).Association(model.ChoosesColumn).Replace(chooses)
		status.Message = model.TagsUpdated
		httpStatus = http.StatusOK
	} else {
		status.Code = model.FailureCode
		status.Message = http.StatusText(http.StatusForbidden)
		httpStatus = http.StatusForbidden
	}
	WriteData(GetJSON(status), httpStatus, w)
}

func getTagInfo(r *http.Request) []string {
	payload := map[string][]string{"Tags": {}}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)
	return payload["Tags"]
}
