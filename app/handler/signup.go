package handler

import (
	"errors"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strings"
)

func SignUp(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	s := status.New()
	s.Message = status.SignupFailure
	credStatus := status.CredentialStatus{}
	// Check if content type is application/json?
	creds, isValidCred := verifyCredentials(r)
	hashedPass, hashErr := Hash(creds.Password)
	if isValidCred && hashErr == nil {
		creds.Password = hashedPass
		user := model.NewUser()
		unameAvailable := !IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, creds.Username, user)
		emailAvailable := !IsSingleRecordActive(db, model.UserTable, model.EmailColumn, creds.Email, user)
		if unameAvailable && emailAvailable {
			err := createUser(db, w, creds, user)
			if err != nil {
				s.Message = status.SignupFailure
				WriteData(GetJSON(s), http.StatusInternalServerError, w)
			} else {
				s.Code = status.SuccessCode
				s.Message = status.SignupSuccess
				WriteData(GetJSON(s), http.StatusCreated, w)
			}
		} else {
			if !unameAvailable {
				credStatus.Username = status.UsernameExists
			}
			if !emailAvailable {
				credStatus.Email = status.EmailExists
			}
			s.Data = credStatus
			WriteData(GetJSON(s), http.StatusConflict, w)
		}
	} else {
		credStatus.Username = status.UsernameAlphaNum
		credStatus.Email = status.ValidEmail
		s.Data = credStatus
		WriteData(GetJSON(s), http.StatusUnprocessableEntity, w)
	}
}

/*
Returning (hash, true) on Hash success otherwise, ("", false) on error.
*/
func Hash(info string) (string, error) {
	// Change cost to 10+ (try to find a way to scale it with hardware?)
	saltedHashPass, err := bcrypt.GenerateFromPassword([]byte(info), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New(status.HashErr)
	}
	return string(saltedHashPass), nil
}

/*
Extracting JSON payload credentials and returning (model, true) if valid, otherwise (model, false).
*/
func verifyCredentials(r *http.Request) (*model.Credentials, bool) {
	c := model.NewCredentials()
	err := extractBody(r, c)
	if err != nil {
		return c, false
	}
	validate := validator.New()
	validate.RegisterValidation("alpha", validateUsername)
	err = validate.Struct(c)
	if err != nil {
		return c, false
	}
	c.Username = strings.ToLower(c.Username)
	c.Email = strings.ToLower(c.Email)
	return c, true
}

/*
Validate username against Regex pattern of being alphanumeric.
*/
func validateUsername(fl validator.FieldLevel) bool {
	matched, _ := regexp.Match("^[a-zA-Z][a-zA-Z0-9_]*$", []byte(fl.Field().String()))
	return matched
}
