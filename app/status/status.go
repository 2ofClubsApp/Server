package status

import (
	"encoding/json"
)

// Status Codes when querying data
const (
	FailureCode = -1
	SuccessCode = 1
)

// Manager addition/removal status messages
const (
	ManagerRemoveFailure   = "unable to remove manager"
	ManagerAdditionFailure = "unable to add manager"
	ManagerRemoveSuccess   = "successfully removed manager"
	ManagerAdditionSuccess = "successfully added manager"
)

// Tag status messages
const (
	TagToggleSuccess = "tag toggled"
	TagCreated       = "successfully created tag"
	TagsCreated      = "successfully created tags"
	TagExists        = "tag already exists"
	TagsUpdated      = "tags updated"
	TagsFound        = "tags found"
	TagNotFound      = "tag doesn't exist"
)

// Authentication status messages
const (
	TokenPairGenerateSuccess = "generated new token pair"
	SignupSuccess            = "signup successful"
	SignupFailure            = "unable to sign up user"
	LoginSuccess             = "successfully logged in"
	LoginFailure             = "username or password is incorrect"
	UsernameExists           = "username already exists"
	UsernameAlphaNum         = "username must start with a letter and can only contain the following characters: a-zA-Z0-9_ and must be 50 characters or less"
	ValidEmail               = "must be a valid email"
	EmailExists              = "email already exists"
	PasswordRequired         = "a password is required"
	LogoutSuccess            = "successfully logged out"
	LogoutFailure            = "unable to logout user"
)

// Password reset status messages
const (
	PasswordUpdateSuccess = "successfully updated password"
	PasswordUpdateFailure = "unable to update password"
	EmailSendFailure      = "unable to send email"
	EmailSendSuccess      = "successfully sent email if user exists"
)

// User & Club status messages
const (
	GetClubManagerSuccess      = "retrieved all club managers"
	ToggleUserSuccess          = "toggled user"
	UserFound                  = "user found"
	UserNotFound               = "user not found"
	ClubFound                  = "club found"
	ClubNotFound               = "club not found"
	ClubCreationSuccess        = "club successfully created"
	ClubCreationFailure        = "unable to create club"
	ClubUpdateSuccess          = "successfully updated club"
	ClubToggleSuccess          = "toggled club"
	ClubUpdateFailure          = "unable to update club"
	GetNonApprovedUsersSuccess = "retrieved all users that require approval"
	GetNonApprovedClubsSuccess = "retrieved all clubs that require approval"
)

// Admin status messages
const (
	ManagerOwnerRequired = "you must be a manager or owner of the club"
	AdminRequired        = "please contact an administrator"
	InvalidTxtFile       = "invalid file: A .txt file is required"
	FileNotFound         = "file doesn't exist"
	UserNotApproved      = "sorry, your account has not been approved yet"
)

// Events status messages
const (
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
)

// Photo upload status messages
const (
	InvalidPhotoFormat = "invalid file: A .jpg or .png file of 10 MB or less is required"
	FileWriteSuccess   = "successfully written to file"
	ClubPhotoNotFound  = "unable to find a photo for the club"
)

// Hashing status messages
const (
	HashErr     = "hashing Error"
	ErrTokenGen = "token generation error"
)

// ErrGeneric represents an generic error that occurred
var ErrGeneric = "an error occurred"

// T - Basic generic struct for any data fetched
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

// CredentialStatus - Listing out requirements for a signup to be successfully created
type CredentialStatus struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// New - Creates a new basic status
// By default, the status code will be a FailureCode (i.e. -1)
func New() *Status {
	return &Status{Code: FailureCode, Message: "", Data: T{}}
}

/*
Display - Returning the JSON representation of a struct.
*/
func (s *Status) Display() string {
	data, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return ""
	}
	return string(data)
}
