package handler

import (
	"errors"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strings"
)

const (
	SignupSuccess = "Signup Successful"
	SignupFailure = "Unable to Sign Up User"
)

func SignUp(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Check if content type is application/json?
	creds, isValid := VerifyCredentials(r)
	hashedPass, err := Hash(creds.Password)
	creds.Password = hashedPass
	status := model.NewStatus()
	status.Message = SignupFailure
	credStatus := model.CredentialStatus{}
	if isValid && (err == nil) {
		user := model.NewUser()
		unameAvailable := !SingleRecordExists(db, model.UserTable, model.UsernameColumn, creds.Username, user)
		emailAvailable := !SingleRecordExists(db, model.UserTable, model.EmailColumn, creds.Email, user)
		if unameAvailable && emailAvailable {
			CreateUser(db, w, creds, user)
		} else {
			if !unameAvailable {
				credStatus.Username = model.UsernameExists
			}
			if !emailAvailable {
				credStatus.Email = model.EmailExists
			}
			status.Data = credStatus
			WriteData(GetJSON(status), http.StatusOK, w)
		}
	} else {
		credStatus.Username = model.UsernameAlphaNum
		credStatus.Email = model.ValidEmail
		status.Data = credStatus
		WriteData(GetJSON(status), http.StatusOK, w)
	}
}

/*
Returning (hash, true) on Hash success otherwise, ("", false) on error.
*/
func Hash(info string) (string, error) {
	// Change cost to 10+ (try to find a way to scale it with hardware?)
	saltedHashPass, err := bcrypt.GenerateFromPassword([]byte(info), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New(model.HashErr)
	}
	return string(saltedHashPass), nil
}

/*
Extracting JSON payload credentials and returning (model, true) if valid, otherwise (model, false).
*/
func VerifyCredentials(r *http.Request) (*model.Credentials, bool) {
	c := model.NewCredentials()
	extractBody(r, c)
	validate := validator.New()
	validate.RegisterValidation("alpha", ValidateUsername)
	err := validate.Struct(c)
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
func ValidateUsername(fl validator.FieldLevel) bool {
	matched, _ := regexp.Match("^[a-zA-Z][a-zA-Z0-9_]*$", []byte(fl.Field().String()))
	return matched
}
