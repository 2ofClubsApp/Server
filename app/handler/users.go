package handler

import (
	"encoding/json"
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/go-playground/validator"
	"github.com/matcornic/hermes/v2"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/*
Create a User given a JSON payload. (See models.User for payload information).
*/
func createUser(db *gorm.DB, w http.ResponseWriter, c *model.Credentials, u *model.User) error {
	u.Credentials = c
	res := db.Create(u)
	if res.Error != nil {
		fmt.Errorf("unable to create user")
	}
	return nil
}

/*
Validating the user request to ensure that they can only access/modify their own data.
True if the requested user has the same username identifier as the token username
*/
func IsValidRequest(username string, r *http.Request) bool {
	claims := GetTokenClaims(r)
	sub := fmt.Sprintf("%v", claims["sub"])
	return sub == username
}

func GetUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	getUserInfo(db, w, r, model.AllUserInfo)
}

/*
Returns all of the Clubs that a User currently manages
*/
func GetUserClubsManage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	getUserInfo(db, w, r, model.AllUserClubsManage)
}

/*
Returns all Events that a User currently attends
*/
func GetUserEventsAttend(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	getUserInfo(db, w, r, model.AllUserEventsAttend)
}

/*
Return partial of all of a users information
Current Supported Information:
	- Users clubs they manage
	- Users events they attend
	- All user info

(See docs for more info and usage)
*/

func getUserInfo(db *gorm.DB, w http.ResponseWriter, r *http.Request, infoType string) {
	var httpStatus int
	var data string
	username := strings.ToLower(getVar(r, model.UsernameVar))
	s := status.New()
	user := model.NewUser()
	if IsValidRequest(username, r) {
		// Defaults will be overridden when obtaining data and being inserted into struct except for null
		userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, username, user)
		if userExists {
			switch strings.ToLower(infoType) {
			case model.AllUserInfo:
				userDisplay := user.Display()
				db.Table(model.UserTable).Preload(model.ManagesColumn).Find(user)
				db.Table(model.UserTable).Preload(model.ChoosesColumn).Find(user)
				db.Table(model.UserTable).Preload(model.AttendsColumn).Find(user)
				userDisplay.Manages = getManages(db, user)
				userDisplay.Tags = filterTags(user.Chooses)
				userDisplay.Attends = user.Attends
				s.Data = userDisplay
			case model.AllUserClubsManage:
				db.Table(model.UserTable).Preload(model.ManagesColumn).Find(user)
				response := make(map[string][]*model.ManagesDisplay)
				response[model.ManagesColumn] = getManages(db, user)
				s.Data = response
			case model.AllUserEventsAttend:
				db.Table(model.UserTable).Preload(model.AttendsColumn).Find(user)
				response := make(map[string][]model.Event)
				response[model.AttendsColumn] = user.Attends
				s.Data = response
			}
			s.Code = status.SuccessCode
			s.Message = status.UserFound
		} else {
			s.Message = status.UserNotFound
		}
		httpStatus = http.StatusOK
	} else {
		s.Message = http.StatusText(http.StatusForbidden)
		httpStatus = http.StatusForbidden
	}
	data = GetJSON(s)
	WriteData(data, httpStatus, w)
}

/*
Extracts the JSON body payload into a given struct (i.e. User, Credentials, etc.)
*/
func extractBody(r *http.Request, s interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(s)
	if err != nil {
		fmt.Errorf("unable to extract JSON payload")
	}
	return nil

}

func getManages(db *gorm.DB, user *model.User) []*model.ManagesDisplay {
	manages := []*model.ManagesDisplay{}
	for _, club := range user.Manages {
		managesDisplay := model.ManagesDisplay{}
		loadClubData(db, &club)
		managesDisplay.Club = club
		managesDisplay.IsOwner = isOwner(db, user, &club)
		manages = append(manages, &managesDisplay)
	}
	return manages
}

/*
Updating the users choice of tags and attended events. Only valid tags will be extracted and added if it's not already.
If an invalid format is provided where there aren't any valid tags to be extracted, the users tag preferences will be reset
*/
func UpdateUserTags(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var httpStatus int
	s := status.New()
	username := strings.ToLower(getVar(r, model.UsernameVar))
	user := model.NewUser()
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, username, user)
	if userExists && IsValidRequest(username, r) {
		httpStatus = http.StatusOK
		// User is guaranteed to have an account (Verified JWT and request is verified)
		chooses := filterTags(extractTags(db, r))
		db.Model(user).Association(model.ChoosesColumn).Replace(chooses)
		s.Code = status.SuccessCode
		s.Message = status.TagsUpdated
	} else {
		s.Message = http.StatusText(http.StatusForbidden)
		httpStatus = http.StatusForbidden
	}
	WriteData(GetJSON(s), httpStatus, w)
}

func UpdateUserPassword(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	newCreds := model.NewPasswordChange()
	extractBody(r, newCreds)
	creds := model.NewCredentials()
	username := strings.ToLower(getVar(r, model.UsernameVar))
	user := model.NewUser()
	s := status.New()
	s.Message = status.PasswordUpdateFailure
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, username, user)
	if userExists {
		if IsValidRequest(user.Username, r) {
			currentPass, ok := getPasswordHash(db, user.Username)
			if bcrypt.CompareHashAndPassword(currentPass, []byte(newCreds.OldPassword)) == nil && ok {
				validate := validator.New()
				creds.Username = user.Username
				creds.Password = newCreds.NewPassword
				creds.Email = user.Email
				validUser := validate.Struct(creds)
				hashedNewPass, hashErr := Hash(newCreds.NewPassword)
				if hashErr == nil && validUser == nil {
					res := db.Model(user).Update(model.PasswordColumn, hashedNewPass)
					if res.Error == nil {
						s.Message = status.PasswordUpdateSuccess
						s.Code = status.SuccessCode
					}
				}
			}
		}
	} else {
		s.Message = status.UserNotFound
	}
	WriteData(GetJSON(s), http.StatusOK, w)
}

func ResetUserPassword(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	creds := model.NewCredentials()
	token := getVar(r, model.TokenVar)
	username := getVar(r, model.UsernameVar)
	user := model.NewUser()
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, username, user)
	s := status.New()
	s.Message = status.PasswordUpdateFailure
	if userExists {
		hash, obtainedHash := getPasswordHash(db, username)
		if obtainedHash {
			if IsValidJWT(token, KF(string(hash))) {
				extractBody(r, creds)
				validate := validator.New()
				// Populate creds struct to validate
				creds.Username = user.Username
				creds.Email = user.Email
				credErr := validate.Struct(creds)
				if newPass, hashErr := Hash(creds.Password); credErr == nil && hashErr == nil {
					res := db.Model(user).Update(model.PasswordColumn, newPass)
					if res.Error == nil {
						s.Code = status.SuccessCode
						s.Message = status.PasswordUpdateSuccess
					}
				}
			}
		}
	}
	WriteData(GetJSON(s), http.StatusOK, w)
}

func RequestResetUserPassword(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	username := strings.ToLower(getVar(r, model.UsernameVar))
	user := model.NewUser()
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, username, user)
	s := status.New()
	s.Message = status.EmailSendFailure
	outputFileName := "template.html"
	if userExists {
		h := hermes.Hermes{
			Product: hermes.Product{
				Name:      os.Getenv("COMPANY_NAME"),
				Link:      os.Getenv("COMPANY_LINK"),
				Logo:      os.Getenv("COMPANY_LOGO"),
				Copyright: os.Getenv("COMPANY_COPYRIGHT"),
			},
		}

		token, jwtErr := GenerateJWT(user.Username, 10, user.Password)
		generateErr := generateEmailTemplate(user, h, outputFileName, token)
		body, fileReadErr := ioutil.ReadFile(outputFileName)
		sendErr := sendEmail(os.Getenv("EMAIL_FROM_HEADER"), user.Email, "Password Reset Request", body)
		if generateErr == nil && fileReadErr == nil && sendErr == nil && jwtErr == nil {
			s.Code = status.SuccessCode
			s.Message = status.EmailSendSuccess
		}
	} else {
		s.Code = status.SuccessCode
		s.Message = status.EmailSendSuccess
	}
	WriteData(GetJSON(s), http.StatusOK, w)
}

func sendEmail(fromEmail string, toEmail string, subject string, body []byte) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", string(body))
	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		return fmt.Errorf("port error")
	} else {
		d := gomail.NewDialer(os.Getenv("EMAIL_HOST"), port, os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"))
		if err := d.DialAndSend(m); err != nil {
			return fmt.Errorf("unable to send email")
		}
	}
	return nil
}
func generateEmailTemplate(user *model.User, h hermes.Hermes, outputFileName string, token string) error {
	email := hermes.Email{
		Body: hermes.Body{
			Intros: []string{"You are receiving this message because you requested to reset your password"},
			Actions: []hermes.Action{
				{
					Instructions: "Click on button below to reset your password:",
					Button: hermes.Button{
						Color: "#DC4D2F",
						Text:  "Reset your password",
						Link:  fmt.Sprintf("http://localhost:8080/resetpassword/%s/%s", user.Username, token),
					},
				},
			},
			Signature: os.Getenv("EMAIL_BODY_SIGNATURE"),
			Outros:    []string{"This link expires in 5 minutes. If you did not request a password reset, please ignore this email."},
			Title:     fmt.Sprintf("Hi %s,", user.Username),
		},
	}
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		return fmt.Errorf("unable to generate email")
	}
	err = ioutil.WriteFile(outputFileName, []byte(emailBody), 0644)
	if err != nil {
		return fmt.Errorf("unable to write email")
	}
	return nil
}
