package status

const (
	// Status Code
	FailureCode = -1
	SuccessCode = 1

	// Manager
	ManagerRemoveFailure   = "Unable to remove manager"
	ManagerAdditionFailure = "Unable to add manager"
	ManagerRemoveSuccess   = "Successfully removed manager"
	ManagerAdditionSuccess = "Successfully added manager"

	// Tags
	TagUpdateError = "Error Updating Tag"
	TagUpdated     = "Tag Updated"
	TagCreated     = "Successfully created tag"
	TagsCreated    = "Successfully created tags"
	TagExists      = "Tag already exists"
	TagsUpdated    = "Tags Updated"
	TagsFound      = "Tags Found"
	TagsGetFailure = "Unable to get tags"
	TagNotFound    = "Tag doesn't exist"

	// Authentication
	SignupSuccess    = "Signup Successful"
	SignupFailure    = "Unable to Sign Up User"
	LoginSuccess     = "Successfully logged in"
	LoginFailure     = "Username or Password is Incorrect"
	UsernameExists   = "Username already exists"
	UsernameAlphaNum = "Username must start with a letter and can only contain the following characters: a-zA-Z0-9_ and must be 50 characters or less"
	ValidEmail       = "Must be a valid email"
	EmailExists      = "Email already exists"

	// Password Update/Reset
	PasswordUpdateSuccess = "Successfully updated password"
	PasswordUpdateFailure = "Unable to update password"
	EmailSendFailure      = "Unable to send email"
	EmailSendSuccess      = "Successfully sent email if user exists"

	// User & Club
	UserUpdated         = "Updated User"
	UserFound           = "User Found"
	UserNotFound        = "User Not Found"
	ClubsFound          = "Clubs Found"
	ClubsNotFound       = "Clubs Not Found"
	ClubFound           = "Club Found"
	ClubNotFound        = "Club Not Found"
	ClubCreationSuccess = "Club successfully created"
	ClubCreationFailure = "Unable to create the Club"
	ClubUpdateSuccess   = "Successfully updated club"
	ClubUpdateFailure   = "Unable to update club"

	// Admin
	GetNonToggledUsersFailure = "Unable to retrieve non toggled users"
	GetNonToggledUsersSuccess = "Retrieved all non toggled users"
	ManagerOwnerRequired      = "You must be a manager or owner of the club"
	AdminRequired             = "Please contact an administrator."
	InvalidTxtFile            = "Invalid File: A .txt file is required"
	UnableToReadFile          = "Unable to read file"
	FileNotFound              = "File doesn't exist"
	UserNotApproved           = "Sorry, your account has not been approved yet"
	ClubNotActive             = "Sorry, this club isn't active yet. Please wait until an administrator activates the club."

	// Events
	CreateEventSuccess         = "Successfully created event"
	CreateEventFailure         = "Unable to create event"
	UpdateEventSuccess         = "Successfully updated event"
	UpdateEventFailure         = "Unable to update event"
	EventNameConstraint        = "Event name must be at least 1 character and a maximum of 50 characters"
	EventDescriptionConstraint = "Event description must be a maximum of 300 characters or less"
	EventLocationConstraint    = "Event location must be a maximum of 100 characters or less"
	EventFeeConstraint         = "Fee must be greater or equal to 0"
	EventFound                 = "Event Found"
	EventNotFound              = "Event not found"
	EventDeleted               = "Event Deleted"
	GetAllEventsFailure        = "Unable to get all events"
	AllEventsFound             = "All Events Found"

	// Photo Upload
	InvalidPhotoFormat   = "Invalid File: A .jpg or .png file of 10 MB or less is required."
	FileCreationFailure  = "Unable to create file"
	FileReadFailure      = "Unable to read file"
	FileWriteFailure     = "Unable to write file"
	FileWriteSuccess     = "Successfully written to file"
	ClubPhotoNotFound = "Unable to find a photo for the club"

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
}

func New() *Status {
	return &Status{Code: FailureCode, Message: "", Data: T{}}
}
