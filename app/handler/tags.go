package handler

import (
	"encoding/json"
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/go-playground/validator"
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
Create a single tag provided the proper JSON request (See the docs for more info)
*/
func CreateTag(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	s := status.New()
	if isAdmin(db, r) {
		payload := map[string]string{}
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&payload)
		tagName := payload["Name"]
		if tagExists(db, tagName) {
			s.Message = status.TagExists
		} else {
			s.Code = status.SuccessCode
			s.Message = status.TagCreated
		}
	} else {
		s.Message = status.AdminRequired
	}
	WriteData(GetJSON(s), http.StatusOK, w)
}

/*
Create tags based on a new line separated list

Refer to docs for file specifications.
*/
func UploadTagsList(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	s := status.New()
	if isAdmin(db, r) {
		file, handler, err := r.FormFile("file")
		if err != nil {
			fmt.Errorf("file doesn't exist: %v", err)
			s.Message = status.FileNotFound
		} else {
			if filepath.Ext(handler.Filename) != ".txt" {
				s.Message = status.InvalidTxtFile
			} else {
				fileContent, err := ioutil.ReadAll(file)
				defer file.Close()
				if err != nil {
					fmt.Errorf("cannot read file: %v", err)
					s.Message = status.UnableToReadFile
				} else {
					for _, tagName := range strings.Split(string(fileContent), "\n") {
						tagName = strings.TrimSpace(tagName)
						tagExists(db, tagName)
					}
					s.Code = status.SuccessCode
					s.Message = status.TagsCreated
				}
			}
		}
	} else {
		s.Message = status.AdminRequired
	}
	WriteData(GetJSON(s), http.StatusOK, w)
}

func GetTags(db *gorm.DB, w http.ResponseWriter, _ *http.Request) {
	s := status.New()
	tags, err := getAllTags(db)
	if err != nil {
		s.Message = status.TagsGetFailure
	} else {
		s.Code = status.SuccessCode
		s.Message = status.TagsFound
		s.Data = tags
	}
	WriteData(GetJSON(s), http.StatusOK, w)
}

func getAllTags(db *gorm.DB) ([]model.Tag, error) {
	var allTags []model.Tag
	result := db.Find(&allTags)
	if result.Error != nil {
		return allTags, fmt.Errorf("unable to get tags")
	}
	return allTags, nil
}

func GetActiveTags(db *gorm.DB, w http.ResponseWriter, _ *http.Request) {
	s := status.New()
	tags, err := getAllTags(db)
	if err != nil {
		s.Message = status.TagsGetFailure
	} else {
		s.Code = status.SuccessCode
		s.Message = status.TagsFound
		s.Data = filterTags(tags)
	}
	WriteData(GetJSON(s), http.StatusOK, w)
}

/*
Extract all tags from payload and returns them as an array of model.Tag
*/
func extractTags(db *gorm.DB, r *http.Request) []model.Tag {
	var chooses []model.Tag
	payload := map[string][]string{"Tags": {}}
	extractBody(r, &payload)
	for _, name := range payload["Tags"] {
		tag := model.NewTag()
		if SingleRecordExists(db, model.TagTable, model.NameColumn, name, tag) {
			chooses = append(chooses, *tag)
		}
	}
	return chooses
}

func ToggleTag(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	s := status.New()
	if isAdmin(db, r) {
		tagName := getVar(r, model.TagNameVar)
		tagName = strings.TrimSpace(tagName)
		tag := model.NewTag()
		if SingleRecordExists(db, model.TagTable, model.NameColumn, tagName, tag) {
			err := db.Model(tag).Update(model.IsActiveColumn, !tag.IsActive).Error
			if err != nil {
				s.Message = status.TagUpdateError
			} else {
				s.Code = status.SuccessCode
				s.Message = status.TagUpdated
			}
		} else {
			s.Message = status.TagNotFound
		}
	} else {
		s.Message = status.AdminRequired
	}
	WriteData(GetJSON(s), http.StatusOK, w)
}

func filterTags(tags []model.Tag) []model.Tag {
	filteredTags := []model.Tag{}
	for _, tag := range tags {
		if tag.IsActive {
			filteredTags = append(filteredTags, tag)
		}
	}
	return filteredTags
}

func flatten(tags []model.Tag) []string {
	flattenedTags := []string{}
	for _, tag := range tags {
		flattenedTags = append(flattenedTags, tag.Name)
	}
	return flattenedTags
}
