package status

import "encoding/json"

const (
	// Status Code
	FailureCode = -1
	SuccessCode = 1

	// Manager
	ManagerRemoveFailure   = "unable to remove manager"
	ManagerAdditionFailure = "unable to add manager"
	ManagerRemoveSuccess   = "successfully removed manager"
	ManagerAdditionSuccess = "successfully added manager"

	// Tags
	TagUpdateError   = "error updating tag"
	TagToggleSuccess = "tag toggled"
	TagUpdated       = "tag updated"
	TagCreated       = "successfully created tag"
	TagsCreated      = "successfully created tags"
	TagExists        = "tag already exists"
	TagsUpdated      = "tags updated"
	TagsFound        = "tags found"
	TagsGetFailure   = "unable to get tags"
	TagNotFound      = "tag doesn't exist"

	// Authentication
	SignupSuccess    = "signup successful"
	SignupFailure    = "unable to sign up user"
	LoginSuccess     = "successfully logged in"
	LoginFailure     = "username or password is incorrect"
	UsernameExists   = "username already exists"
	UsernameAlphaNum = "username must start with a letter and can only contain the following characters: a-zA-Z0-9_ and must be 50 characters or less"
	ValidEmail       = "must be a valid email"
	EmailExists      = "email already exists"
	PasswordRequired = "a password is required"

	// Password Update/Reset
	PasswordUpdateSuccess = "successfully updated password"
	PasswordUpdateFailure = "unable to update password"
	EmailSendFailure      = "unable to send email"
	EmailSendSuccess      = "successfully sent email if user exists"

	// User & Club
	ToggleUserSuccess   = "toggled user"
	UserUpdated         = "updated user"
	UserFound           = "user found"
	UserNotFound        = "user not found"
	ClubsFound          = "clubs found"
	ClubsNotFound       = "clubs not found"
	ClubFound           = "club found"
	ClubNotFound        = "club not found"
	ClubCreationSuccess = "club successfully created"
	ClubCreationFailure = "unable to create club"
	ClubUpdateSuccess   = "successfully updated club"
	ClubToggleSuccess   = "toggled club"
	ClubUpdateFailure   = "unable to update club"

	// Admin
	GetNonToggledUsersFailure = "unable to retrieve non toggled users"
	GetNonToggledUsersSuccess = "retrieved all non toggled users"
	ManagerOwnerRequired      = "you must be a manager or owner of the club"
	AdminRequired             = "please contact an administrator"
	InvalidTxtFile            = "invalid file: A .txt file is required"
	UnableToReadFile          = "unable to read file"
	FileNotFound              = "file doesn't exist"
	UserNotApproved           = "sorry, your account has not been approved yet"

	// Events
	CreateEventSuccess         = "successfully created event"
	CreateEventFailure         = "unable to create event"
	UpdateEventSuccess         = "successfully updated event"
	UpdateEventFailure         = "unable to update event"
	EventNameConstraint        = "event name must be at least 1 character and a maximum of 50 characters"
	EventDescriptionConstraint = "event description must be a maximum of 300 characters or less"
	EventLocationConstraint    = "event location must be a maximum of 100 characters or less"
	EventFeeConstraint         = "fee must be greater or equal to 0"
	EventFound                 = "event found"
	EventNotFound              = "event not found"
	EventDeleted               = "event deleted"
	GetAllEventsFailure        = "unable to get all events"
	AllEventsFound             = "all events found"
	EventUnattendSuccess       = "event unattended"
	EventAttendSuccess         = "event attended"

	// Photo Upload
	InvalidPhotoFormat  = "invalid file: A .jpg or .png file of 10 MB or less is required"
	FileCreationFailure = "unable to create file"
	FileReadFailure     = "unable to read file"
	FileWriteFailure    = "unable to write file"
	FileWriteSuccess    = "successfully written to file"
	ClubPhotoNotFound   = "unable to find a photo for the club"

	// Hashing
	HashErr     = "hashing Error"
	ErrTokenGen = "token generation error"

	ErrGeneric = "an error occurred"
)

//
type T struct{}

// Status Struct used as the standard response when querying the API
type Status struct {
	// Codes will either be returned as 1 or -1
	// 1 - Success
	// -1 - Failure
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Listing out requirements for a signup to be successfully created
type CredentialStatus struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New() *Status {
	return &Status{Code: FailureCode, Message: "", Data: T{}}
}

/*
Returning the JSON representation of a struct.
*/
func (s *Status) Display() string {
	data, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return ""
	}
	return string(data)
}
