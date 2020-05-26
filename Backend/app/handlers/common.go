package handlers

import (
	"../models"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
)

/*
	Common methods shared amongst the different models
*/

func RecordExists(db *gorm.DB, column string, val string, t interface{}) bool {
	return !db.Where(column+"= ?", val).First(t).RecordNotFound()
}

func ExtractPersonInfo(r *http.Request) models.Person {
	decoder := json.NewDecoder(r.Body)
	p := models.NewPerson()
	decoder.Decode(&p)
	return p
}

func ParseJSON(response interface{}) string {
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

func WriteData(data string, code int, w http.ResponseWriter) int {
	w.WriteHeader(code)
	n, err := fmt.Fprint(w, string(data))
	if err != nil {
		return -1
	}
	return n
}
