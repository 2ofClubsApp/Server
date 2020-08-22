package handler

import (
	"encoding/json"
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

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
			return db.Create(tag).Error != nil
		}
		return false
	}
	return true
}

// CreateTag - Create a single tag provided the proper JSON request (See the docs for more info)
func CreateTag(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	if isAdmin(db, r) {
		payload := map[string]string{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&payload)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
		}
		tagName := payload["Name"]
		if tagExists(db, tagName) {
			s.Message = status.TagExists
			return http.StatusConflict, nil
		}
		s.Code = status.SuccessCode
		s.Message = status.TagCreated
		return http.StatusCreated, nil
	}
	s.Message = status.AdminRequired
	return http.StatusForbidden, nil
}

// UploadTagsList - Create tags based on a new line separated list
// Refer to docs for file specifications.
func UploadTagsList(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	if isAdmin(db, r) {
		file, handler, err := r.FormFile("file")
		if err != nil {
			s.Message = status.FileNotFound
			return http.StatusBadRequest, nil
		}
		if filepath.Ext(handler.Filename) != ".txt" {
			s.Message = status.InvalidTxtFile
			return http.StatusUnsupportedMediaType, nil
		}
		fileContent, err := ioutil.ReadAll(file)
		defer file.Close()
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
		}
		for _, tagName := range strings.Split(string(fileContent), "\n") {
			tagName = strings.TrimSpace(tagName)
			tagExists(db, tagName)
		}
		s.Code = status.SuccessCode
		s.Message = status.TagsCreated
		return http.StatusCreated, nil
	}
	s.Message = status.AdminRequired
	return http.StatusForbidden, nil

}

// GetTags - Obtaining all tags (both active and inactive)
func GetTags(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, _ *http.Request, s *status.Status) (int, error) {
	tags, err := getAllTags(db)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
	}
	s.Code = status.SuccessCode
	s.Message = status.TagsFound
	s.Data = tags
	return http.StatusOK, nil
}

// Helper function to obtain all tags
func getAllTags(db *gorm.DB) ([]model.Tag, error) {
	var allTags []model.Tag
	result := db.Find(&allTags)
	if result.Error != nil {
		return allTags, fmt.Errorf("unable to get tags")
	}
	return allTags, nil
}

// GetActiveTags - Obtaining all active tags
func GetActiveTags(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, _ *http.Request, s *status.Status) (int, error) {
	tags, err := getAllTags(db)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
	}
	s.Code = status.SuccessCode
	s.Message = status.TagsFound
	s.Data = filterTags(tags)
	return http.StatusOK, nil
}

// Extract all tags from payload and returns them as an array of model.Tag
func extractTags(db *gorm.DB, r *http.Request) []model.Tag {
	var chooses []model.Tag
	payload := map[string][]string{"Tags": {}}
	extractBody(r, &payload) // Error check here later
	for _, name := range payload["Tags"] {
		tag := model.NewTag()
		if SingleRecordExists(db, model.TagTable, model.NameColumn, name, tag) {
			chooses = append(chooses, *tag)
		}
	}
	return chooses
}

// ToggleTag - Toggling tags as either active or inactive
func ToggleTag(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	if isAdmin(db, r) {
		tagName := getVar(r, model.TagNameVar)
		tagName = strings.TrimSpace(tagName)
		tag := model.NewTag()
		if SingleRecordExists(db, model.TagTable, model.NameColumn, tagName, tag) {
			err := db.Model(tag).Update(model.IsActiveColumn, !tag.IsActive).Error
			if err != nil {
				return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
			}
			s.Code = status.SuccessCode
			s.Message = status.TagToggleSuccess
			return http.StatusOK, nil
		}
		s.Message = status.TagNotFound
		return http.StatusNotFound, nil
	}
	s.Message = status.AdminRequired
	return http.StatusForbidden, nil
}

// Filtering and returning []model.Tag that are active
func filterTags(tags []model.Tag) []model.Tag {
	filteredTags := []model.Tag{}
	for _, tag := range tags {
		if tag.IsActive {
			filteredTags = append(filteredTags, tag)
		}
	}
	return filteredTags
}

// Flatten []model.Tag and return a list of tag names
func flatten(tags []model.Tag) []string {
	flattenedTags := []string{}
	for _, tag := range tags {
		flattenedTags = append(flattenedTags, tag.Name)
	}
	return flattenedTags
}
