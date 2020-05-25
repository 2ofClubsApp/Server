package handlers

import (
	"../../models"
	"../Status"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

func extractStudentInfo(r *http.Request) models.Student {
	// Need to lowercase username
	d := json.NewDecoder(r.Body)
	p := models.Person{}
	s := models.Student{}
	d.Decode(&p)
	s.Person = p
	fmt.Println("Create Student")
	fmt.Println(s)
	return s
}
func CreateStudent(db *gorm.DB, _ http.ResponseWriter, r *http.Request) {
	// Check if content type is application/json
	s := extractStudentInfo(r)
	db.Create(&s)

}

func GetStudent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars[models.ColumnUsername]
	fmt.Printf("Getting student %s\n", username)
	// Defaults will be overridden when inserting data into struct
	s := models.Student{Tags: []models.Tag{}, Attends: []models.Event{}, Swipes: []models.Club{}}
	isFound := recordExists(db, models.ColumnUsername, username, &s)
	fmt.Println(isFound)
	ss := Status.Status{}
	if isFound {
		fmt.Println("User Not found")
		ss.Status = Status.FAILURE
		ss.Data = s
		fmt.Fprintf(w, "{Message: Student not found}")
		return
	} else {
		ss.Status = Status.SUCCESS
		ss.Data = s
		data, err := json.Marshal(ss)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(data))
	}

}

func UpdateStudent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Student")
}

func recordExists(db *gorm.DB, a string, b string, s *models.Student) bool {
	return db.Where(a+"= ?", b).First(s).RecordNotFound()
}
