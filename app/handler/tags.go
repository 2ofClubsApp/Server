package handler

import (
	"../model"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

// Note: Tags such as Computer Science and ComputerScience are different, should we account for this or is this a user fault?
/*
Returns true if the tag already exists in the database, false otherwise.
If false, the tag will be created and inserted into the database.
*/
func tagExists(db *gorm.DB, tagName string) bool {
	validate := validator.New()
	if !SingleRecordExists(db, model.TagTable, model.NameColumn, tagName, model.NewTag()) {
		tag := model.NewTag()
		tag.Name = tagName
		tag.IsActive = true
		if validate.Struct(tag) == nil {
			db.Create(tag)
		}
		return false
	}
	return true
}

/*
Create a tag based on the name provided by the request URL
*/
func CreateTag(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	status := model.NewStatus()
	if isAdmin(db, r) {
		vars := mux.Vars(r)
		tagName := vars["tag"]
		tagName = strings.TrimSpace(tagName)
		if tagExists(db, tagName) {
			status.Message = model.TagExists
			status.Code = model.FailureCode
		} else {
			status.Message = model.TagCreated
		}
	} else {
		status.Code = model.FailureCode
		status.Message = model.AdminRequired
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}

/*
Create tags based on a new line separated list

Refer to docs for file specifications.
*/
func UploadTagsList(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	status := model.NewStatus()
	if isAdmin(db, r) {
		file, handler, err := r.FormFile("file")
		if err != nil {
			fmt.Errorf("file doesn't exist: %v", err)
			return
		}
		if filepath.Ext(handler.Filename) != ".txt" {
			status.Code = model.FailureCode
			status.Message = model.InvalidFile
		} else {
			fileContent, err := ioutil.ReadAll(file)
			defer file.Close()
			if err != nil {
				fmt.Errorf("cannot read file: %v", err)
				return
			}
			for _, tagName := range strings.Split(string(fileContent), "\n") {
				tagName = strings.TrimSpace(tagName)
				tagExists(db, tagName)
			}
			status.Message = model.TagsCreated
		}
	} else {
		status.Code = model.FailureCode
		status.Message = model.AdminRequired
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}

func GetTags(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	status := model.NewStatus()
	var allTags []model.Tag
	var tagsList []string
	type TagData struct {
		Tags []string
	}
	result := db.Find(&allTags)
	if result.Error != nil {
		fmt.Errorf("unable to get tags: %v", result.Error)
		return
	}
	for _, tag := range allTags {
		tagsList = append(tagsList, tag.Name)
	}
	status.Message = model.TagsFound
	status.Data = TagData{Tags: tagsList}
	WriteData(GetJSON(status), http.StatusOK, w)
}

/*
Extract all tags from payload and returns them as an array of model.Tag
*/
func extractTags(db *gorm.DB, r *http.Request) []model.Tag {
	var chooses []model.Tag
	for _, name := range getTagInfo(r) {
		tag := model.NewTag()
		if SingleRecordExists(db, model.TagTable, model.NameColumn, name, tag) {
			tag.Name = name
			chooses = append(chooses, *tag)
		}
	}
	return chooses
}

func DeleteTag(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	status := model.NewStatus()
	if isAdmin(db, r) {
		vars := mux.Vars(r)
		tagName := vars["tag"]
		tagName = strings.TrimSpace(tagName)
		tag := model.NewTag()
		if SingleRecordExists(db, model.TagTable, model.NameColumn, tagName, tag) {
			db.Delete(tag)
			status.Message = model.TagDelete
		} else {
			status.Code = -1
			status.Message = model.TagNotFound
		}
	} else {
		status.Code = -1
		status.Message = model.AdminRequired
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}

func flatten(tags []model.Tag) []string {
	flattenTags := []string{}
	for _, tag := range tags {
		flattenTags = append(flattenTags, tag.Name)
	}
	return flattenTags
}
