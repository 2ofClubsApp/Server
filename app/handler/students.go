package handler

import (
	"../model"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

/*
Create student given a student JSON payload. (See models.Student for payload information).
 */
func CreateStudent(db *gorm.DB, w http.ResponseWriter, u *model.User, s *model.Student) {
	pass, isHashed := Hash(u.Password)
	u.Password = pass
	s.User = u
	found := RecordExists(db, model.StudentTable, model.UsernameColumn, s.Username, s)
	status := model.NewStatus()
	if !found && isHashed {
		db.Create(&s)
		WriteData(GetJSON(status), http.StatusOK, w)
	} else {
		status.Message = ErrSignUp
		WriteData(GetJSON(status), http.StatusOK, w)
	}
}

func GetStudent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	data := ""
	vars := mux.Vars(r)
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
		data = GetJSON(ss)
	} else {
		data = http.StatusText(http.StatusForbidden)
		status = http.StatusForbidden
	}
	WriteData(data, status, w)

}

func UpdateStudent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Student")
}
