package handler

import (
	"../model"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

func CreateStudent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Check if content type is application/json?
	s := model.NewStudent()
	p := ExtractPersonInfo(r)
	pass, isHashed := Hash(p.Password)
	username, isHashed := Hash(p.Username)
	email, isHashed := Hash(p.Email)
	p.Password = pass
	p.Username = username
	p.Email = email
	s.Person = p
	found := RecordExists(db, "student", model.UsernameColumn, s.Username, s)
	if !found && isHashed {
		if tp, err := GetTokenPair(s.Username, 5, 60*24); err == nil {
			db.Create(&s)
			c := GenerateCookie(model.RefreshToken, tp.RefreshToken)
			http.SetCookie(w, c)
			WriteData(tp.AccessToken, http.StatusOK, w)
		}
	}
}

func GetStudent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
<<<<<<< Updated upstream
	username := vars[model.ColumnUsername]
	s := model.NewStudent()
	ss := model.NewStatus()
	// Defaults will be overridden when obtaining data and being inserted into struct except for null
	found := RecordExists(db, model.ColumnUsername, username, s)
	if !found {
		ss.Message = model.FAILURE
=======
	username := vars[model.UsernameColumn]
	if ValidateUserReq(username, r) {
		s := model.NewStudent()
		ss := model.NewStatus()
		// Defaults will be overridden when obtaining data and being inserted into struct except for null
		found := RecordExists(db, "student", model.UsernameColumn, username, s)
		if !found {
			ss.Message = model.Failure
		} else {
			ss.Data = s
		}
		data = ParseJSON(ss)
>>>>>>> Stashed changes
	} else {
		ss.Data = s
	}
	WriteData(ParseJSON(ss), http.StatusOK, w)
}

func UpdateStudent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Student")
}
