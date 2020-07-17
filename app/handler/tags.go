package handler

import (
	"../model"
	"bufio"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
)

func CreateTags(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	for _, name := range getTagInfo(r){
		if !SingleRecordExists(db, model.TagTable, model.NameColumn, name, model.NewTag()){
			tag := model.NewTag()
			tag.Name = name
			db.Create(tag)
		} else {
			fmt.Println("Record already exists")
		}
	}
}

func getTagInfo(r *http.Request) []string {
	payload := map[string][]string{"Tags": []string{}}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)
	return payload["Tags"]

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
