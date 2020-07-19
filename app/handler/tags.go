package handler

import (
	"../model"
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Note: Tags such as Computer Science and ComputerScience are different, should we account for this or is this a user fault?
func CreateTag(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	status := model.NewStatus()
	vars := mux.Vars(r)
	tagName := vars["tag"]
	tagName = strings.TrimSpace(tagName)
	if !SingleRecordExists(db, model.TagTable, model.NameColumn, tagName, model.NewTag()) {
		tag := model.NewTag()
		tag.Name = tagName
		db.Create(tag)
		status.Message = model.TagCreated
	} else {
		status.Message = model.TagFound
		status.Code = model.FailureCode
	}
	WriteData(GetJSON(status), http.StatusOK, w)

}

func UploadTagsList(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("No file provided")
		return
	}
	defer file.Close()
	f, err := os.OpenFile(handler.Filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error Opening File")
		return
	}
	defer f.Close()
	_, err1 := io.Copy(f, file)
	fc, _ := ioutil.ReadAll(file)
	fmt.Println(fc)
	if err1 != nil {
		fmt.Println("Error copying contents from other file")
		return
	}
	scanner := bufio.NewScanner(f)
	a := scanner.Scan()
	fmt.Println(a)
	for a {
		fmt.Println(scanner.Text())
	}
}

func GetTags(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Club Tags")
}
