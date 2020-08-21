package handler

import (
	"errors"
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strings"
)

// SignUp will create a user given valid credentials
// See model.Credentials or docs for username and email constraints
func SignUp(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	s.Message = status.SignupFailure
	credStatus := status.CredentialStatus{}
	creds, isValidCred := verifyCredentials(r)
	hashedPass, hashErr := Hash(creds.Password)
	if isValidCred == nil && hashErr == nil {
		creds.Password = hashedPass
		user := model.NewUser()
		unameAvailable := !SingleRecordExists(db, model.UserTable, model.UsernameColumn, creds.Username, user)
		emailAvailable := !SingleRecordExists(db, model.UserTable, model.EmailColumn, creds.Email, user)
		if unameAvailable && emailAvailable {
			err := createUser(db, creds, user)
			if err != nil {
				return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
			}
			s.Code = status.SuccessCode
			s.Message = status.SignupSuccess
			return http.StatusCreated, nil
		}
		if !unameAvailable {
			credStatus.Username = status.UsernameExists
		}
		if !emailAvailable {
			credStatus.Email = status.EmailExists
		}
		s.Data = credStatus
		return http.StatusConflict, nil
	}
	credStatus.Username = status.UsernameAlphaNum
	credStatus.Email = status.ValidEmail
	credStatus.Password = status.PasswordRequired
	s.Data = credStatus
	return http.StatusUnprocessableEntity, nil
}

// Hash - Returning (hash, true) on Hash success otherwise, ("", false) on error.
func Hash(info string) (string, error) {
	// Change cost to 10+ (try to find a way to scale it with hardware?)
	saltedHashPass, err := bcrypt.GenerateFromPassword([]byte(info), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New(status.HashErr)
	}
	return string(saltedHashPass), nil
}

// Extracting JSON payload credentials and returning (model, true) if valid, otherwise (model, false).
func verifyCredentials(r *http.Request) (*model.Credentials, error) {
	c := model.NewCredentials()
	err := extractBody(r, c)
	if err != nil {
		return c, fmt.Errorf("unable to extract credentials")
	}
	validate := validator.New()
	if err := validate.RegisterValidation("alpha", validateUsername); err != nil {
		return c, fmt.Errorf("struct validation error")
	}
	err = validate.Struct(c)
	if err != nil {
		return c, fmt.Errorf("struct validation error")
	}
	c.Username = strings.ToLower(c.Username)
	c.Email = strings.ToLower(c.Email)
	return c, nil
}

/*
Validate username against Regex pattern of being alphanumeric.
*/
func validateUsername(fl validator.FieldLevel) bool {
	matched, _ := regexp.Match("^[a-zA-Z][a-zA-Z0-9_]*$", []byte(fl.Field().String()))
	return matched
}
