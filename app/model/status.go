package model

const (
	ManagerRemoveFailure   = "Unable to remove manager"
	ManagerAdditionFailure = "Unable to add manager"
	ManagerRemoveSuccess   = "Successfully removed manager"
	ManagerAdditionSuccess = "Successfully added manager"

	TagUpdateError = "Error Updating Tag"
	TagUpdated     = "Tag Updated"
	TagCreated     = "Successfully created tag"
	TagsCreated    = "Successfully created tags"
	TagExists      = "Tag already exists"
	TagsUpdated    = "Tags Updated"
	TagsFound      = "Tags Found"
	TagsGetFailure = "Unable to get tags"
	TagNotFound    = "Tag doesn't exist"

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

	UsernameExists   = "Username already exists"
	UsernameAlphaNum = "Username must start with a letter and can only contain the following characters: a-zA-Z0-9_ and must be 50 characters or less"
	ValidEmail       = "Must be a valid email"
	EmailExists      = "Email already exists"

	GetNonToggledUsersFailure = "Unable to retrieve non toggled users"
	GetNonToggledUsersSuccess = "Retrieved all non toggled users"

	ManagerOwnerRequired = "You must be a manager or owner of the club"
	AdminRequired        = "Please contact an administrator."
	InvalidFile          = "Invalid File: A .txt file is required"
	UnableToReadFile     = "Unable to read file"
	FileNotFound         = "File doesn't exist"

	UserNotApproved       = "Sorry, your account has not been approved yet"
	ClubNotActive         = "Sorry, this club isn't active yet. Please wait until an administrator activates the club."
	LoginSuccess          = "Successfully logged in"
	LoginFailure          = "Username or Password is Incorrect"
	PasswordUpdateSuccess = "Successfully updated password"
	PasswordUpdateFailure = "Unable to update password"

	CreateEventSuccess = "Successfully created event"
	CreateEventFailure = "Unable to create event"
	UpdateEventSuccess = "Successfully updated event"
	UpdateEventFailure = "Unable to update event"

	EventNameConstraint        = "Event name must be at least 1 character and a maximum of 50 characters"
	EventDescriptionConstraint = "Event description must be a maximum of 300 characters or less"
	EventLocationConstraint    = "Event location must be a maximum of 100 characters or less"
	EventFeeConstraint         = "Fee must be greater or equal to 0"
	EventFound                 = "Event Found"
	EventNotFound              = "Event not found"
	EventDeleted               = "Event Deleted"
	GetAllEventsFailure        = "Unable to get all events"
	AllEventsFound             = "All Events Found"

	HashErr     = "hashing Error"
	ErrTokenGen = "token generation error"

	EmailSendFailure = "Unable to send email"
	EmailSendSuccess = "Successfully sent email if user exists"

	ErrGeneric  = "an error occurred"
	FailureCode = -1
	SuccessCode = 1
)

type T struct{}

type Status struct {
	Code    int
	Message string
	Data    interface{}
}

type CredentialStatus struct {
	Username string
	Email    string
}

func NewStatus() *Status {
	return &Status{Code: FailureCode, Message: "", Data: T{}}
}
