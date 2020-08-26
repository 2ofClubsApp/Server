package handler

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GetClubs returns all of the clubs that are
func GetClubs(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	uname := getVar(r, model.UsernameVar)
	user := model.NewUser()
	unameExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	claims := GetTokenClaims(ExtractToken(r))
	tokenUname := fmt.Sprintf("%v", claims["sub"])
	if uname != tokenUname {
		s.Message = http.StatusText(http.StatusForbidden)
		return http.StatusForbidden, nil
	}
	if !unameExists {
		s.Message = status.UserNotFound
		return http.StatusNotFound, nil
	}
	activeTags := flatten(filterTags(extractTags(db, r)))
	allClubs := []model.Club{}
	clubsWithTag := []model.Club{}
	clubsWithTagNonFavourited := []model.Club{}
	if db.Where(model.ActiveColumn+"= ?", true).Find(&allClubs).Error != nil {
		return http.StatusInternalServerError, fmt.Errorf("unable to obtain all active clubs")
	}
	if len(activeTags) != 0 {
		for _, c := range allClubs {
			loadClubData(db, &c)
			set := false
			for _, tag := range c.Sets {
				if tagInSlice(tag.Name, activeTags) && !set {
					clubsWithTag = append(clubsWithTag, c)
					set = true
				}
			}
		}
	} else {
		for _, c := range allClubs {
			loadClubData(db, &c)
			clubsWithTag = append(clubsWithTag, c)
		}
	}
	res := db.Table(model.UserTable).Preload(model.SwipedColumn).Find(user)
	if res.Error != nil {
		return http.StatusInternalServerError, fmt.Errorf("unable to load user swiped clubs")
	}
	for _, filteredTagClub := range clubsWithTag {
		swiped := false
		for _, club := range user.Swiped {
			if club.ID == filteredTagClub.ID {
				swiped = true
				break
			}
		}
		if !swiped {
			clubsWithTagNonFavourited = append(clubsWithTagNonFavourited, filteredTagClub)
		}
	}
	s.Code = status.SuccessCode
	s.Message = status.GetFilteredNonSwipedClubsSuccess
	s.Data = clubsWithTagNonFavourited
	return http.StatusOK, nil
}

// Returns true if the tagName is in the slice, false otherwise
// Helper function for GetClubs
func tagInSlice(tag string, activeTags []string) bool {
	for _, tagName := range activeTags {
		if tag == tagName {
			return true
		}
	}
	return false
}

// UpdateClub - Update club with new information
// See model.Club or docs for club attribute specifications
func UpdateClub(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	user := model.NewUser()
	claims := GetTokenClaims(ExtractToken(r))
	uname := fmt.Sprintf("%v", claims["sub"])
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	if !userExists {
		s.Message = status.UserNotFound
		return http.StatusNotFound, nil
	}
	if !isManager(db, user, club) {
		s.Code = status.FailureCode
		s.Message = http.StatusText(http.StatusForbidden)
		return http.StatusForbidden, nil
	}
	updatedClub := model.NewClubUpdate()
	extractBody(r, updatedClub)
	validate := validator.New()
	err := validate.Struct(updatedClub)
	if err != nil {
		s.Message = status.ClubUpdateFailure
		return http.StatusUnprocessableEntity, nil
	}
	club.Bio = updatedClub.Bio
	club.Size = updatedClub.Size
	res := db.Model(club).Select(model.BioColumn, model.SizeColumn).Updates(club)
	if res.Error != nil {
		return http.StatusInternalServerError, fmt.Errorf("unable to update club")
	}
	s.Code = status.SuccessCode
	s.Message = status.ClubUpdateSuccess
	return http.StatusOK, nil
}

// CreateClub - Creating a club (You must have an active user account first)
// See model.Club or the docs for club information constraints
func CreateClub(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	claims := GetTokenClaims(ExtractToken(r))
	user := model.NewUser()
	uname := fmt.Sprintf("%v", claims["sub"])
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	club := model.NewClub()
	if extractBody(r, club) != nil {
		return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
	}
	validate := validator.New()
	err := validate.Struct(club)
	if err != nil {
		return http.StatusUnprocessableEntity, nil
	}
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.NameColumn, club.Name, model.NewClub())
	emailExists := SingleRecordExists(db, model.ClubTable, model.EmailColumn, club.Email, model.NewClub())
	if !emailExists && !clubExists && userExists && err == nil {
		// Clearing any links supplied
		if club.Logo != "" {
			club.Logo = ""
		}
		err := db.Model(user).Association(model.ManagesColumn).Append(club)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
		}
		res := db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ? AND is_owner = ?", user.ID, club.ID, false).Update(model.IsOwnerColumn, true)
		if res.Error != nil {
			return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
		}
		s.Message = status.ClubCreationSuccess
		s.Code = status.SuccessCode
		return http.StatusCreated, nil
	}
	s.Message = status.ClubCreationFailure
	return http.StatusUnprocessableEntity, nil
}

// GetClubPhoto - Obtaining a club profile photo (if it exists)
func GetClubPhoto(db *gorm.DB, _ *redis.Client, w http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	clubID := getVar(r, model.ClubIDVar)
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, model.NewClub())
	if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	_, err := ioutil.ReadDir("./images")
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("dir doesn't exist")
	}
	path := fmt.Sprintf("images/%s.png", clubID)
	if _, err := os.Stat(path); err != nil {
		s.Message = status.ClubPhotoNotFound
		return http.StatusNotFound, nil
	}
	img, err := os.Open(path)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("unable to open photo")
	}
	defer img.Close()
	//w.Header().Set("Content-Type", "image/jpeg")
	_, err = io.Copy(w, img)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("unable to read file contents")
	}
	s.Code = status.SuccessCode
	return http.StatusOK, nil
}

// UploadClubPhoto - Uploading a club photo
// Club photo file size upload is 10 MB max
func UploadClubPhoto(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	maxMem := int64(10 << 20)
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	user := model.NewUser()
	claims := GetTokenClaims(ExtractToken(r))
	username := fmt.Sprintf("%v", claims["sub"])
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, username, user)
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	if !(clubExists && userExists && isManager(db, user, club)) {
		s.Message = http.StatusText(http.StatusForbidden)
		return http.StatusForbidden, nil
	}
	// Max 10MB upload file
	err := r.ParseMultipartForm(maxMem) // 2^20
	if err != nil {
		s.Message = status.InvalidPhotoSize
		return http.StatusBadRequest, nil
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		s.Message = status.FileNotFound
		return http.StatusBadRequest, nil
	}
	if filepath.Ext(handler.Filename) != ".png" && filepath.Ext(handler.Filename) != ".jpg" {
		s.Message = status.InvalidPhotoFormat
		return http.StatusUnsupportedMediaType, nil
	}
	defer file.Close()
	fileName := fmt.Sprintf("images/%v.png", club.ID)
	tempFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("unable to create temp file")
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("unable to read file contents")
	}
	_, err = tempFile.Write(fileBytes)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("unable to write filebytes")
	}
	if club.Logo == "" {
		club.Logo = fmt.Sprintf("/photos/clubs/%v", club.ID)
		if db.Model(club).Select(model.LogoColumn).Updates(club).Error != nil {
			return http.StatusInternalServerError, fmt.Errorf("unable to set club logo path")
		}
	}
	s.Code = status.SuccessCode
	s.Message = status.FileUploadSuccess
	return http.StatusOK, nil

}

// UpdateClubTags - Updating user club tags
// All old tags will be overrided with the new set of tags provided
func UpdateClubTags(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	claims := GetTokenClaims(ExtractToken(r))
	username := fmt.Sprintf("%v", claims["sub"])
	user := model.NewUser()
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, username, user)
	if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	// Must check with both user and club existing in the event that a user gets deleted but you manage to get a hold of their access token
	if !(userExists && clubExists && isManager(db, user, club)) {
		s.Message = http.StatusText(http.StatusForbidden)
		return http.StatusForbidden, nil
	}
	tags := filterTags(extractTags(db, r))
	err := db.Model(club).Association(model.SetsColumn).Replace(tags)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("unable to get club tags")
	}
	s.Message = status.TagsUpdated
	s.Code = status.SuccessCode
	return http.StatusOK, nil
}

// GetClub - Obtaining all information about a club
func GetClub(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	return getClubInfo(db, r, model.AllClubInfo, s)
}

// A helper function returning specific or all parts of a club
func getClubInfo(db *gorm.DB, r *http.Request, infoType string, s *status.Status) (int, error) {
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	switch strings.ToLower(infoType) {
	case model.AllClubInfo:
		if loadClubData(db, club) != nil {
			return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
		}
		s.Data = club
		s.Message = status.ClubFound
	case model.AllClubEventsHost:
		clubEvents := make(map[string][]model.Event)
		res := db.Table(model.ClubTable).Preload(model.HostsColumn).Find(club)
		if res.Error != nil {
			return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
		}
		clubEvents[strings.ToLower(model.HostsColumn)] = club.Hosts
		s.Data = clubEvents
		s.Message = status.ClubEventFound
	}
	s.Code = status.SuccessCode
	return http.StatusOK, nil
}

// GetClubEvents - Obtaining all events that a club hosts
func GetClubEvents(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	return getClubInfo(db, r, model.AllClubEventsHost, s)
}

// Returning club relational data (i.e. Club & Tags and Club & Events) and populating to Club struct
func loadClubData(db *gorm.DB, club *model.Club) error {
	if db.Table(model.ClubTable).Preload(model.SetsColumn).Find(club).Error != nil {
		return fmt.Errorf("unable to obtain club tags")
	}
	if db.Table(model.ClubTable).Preload(model.HostsColumn).Find(club).Error != nil {
		return fmt.Errorf("unable to obtain club events hosted")
	}
	club.Sets = filterTags(club.Sets)
	if len(club.Hosts) == 0 {
		club.Hosts = []model.Event{}
	}
	return nil
}

/*
Returns true if the user is an owner of the club, false otherwise
*/
func isOwner(db *gorm.DB, user *model.User, club *model.Club) bool {
	userClub := model.NewUserClub()
	db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ?", user.ID, club.ID).First(userClub)
	return userClub.IsOwner
}

// Returns true if the user is a manager of a club, false otherwise
func isManager(db *gorm.DB, user *model.User, club *model.Club) bool {
	userClub := model.NewUserClub()
	res := db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ?", user.ID, club.ID).First(userClub)
	return res.Error == nil
}

// RemoveManager - Removing a manager from a club
// Note: The user removing the manager must be a club owner
func RemoveManager(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	return editManagers(db, r, model.OpRemove, s)
}

// AddManager - Adding a manager to a club
// Note: The user adding the manager must be a club owner
func AddManager(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	return editManagers(db, r, model.OpAdd, s)
}

/*
Adding or removing managers and their associations to a particular club
*/
func editManagers(db *gorm.DB, r *http.Request, op string, s *status.Status) (int, error) {
	// Default messages set to manager addition, otherwise manager removal
	var successMessage = status.ManagerAdditionSuccess
	var failureMessage = status.ManagerAdditionFailure
	if op == model.OpRemove {
		successMessage = status.ManagerRemoveSuccess
		failureMessage = status.ManagerRemoveFailure
	}
	claims := GetTokenClaims(ExtractToken(r))
	clubOwnerUsername := fmt.Sprintf("%v", claims["sub"])
	newManagerUname := strings.ToLower(getVar(r, model.UsernameVar))
	clubID := getVar(r, model.ClubIDVar)
	owner := model.NewUser()
	newManager := model.NewUser()
	club := model.NewClub()
	// Added user must exist
	ownerExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, clubOwnerUsername, owner)
	// If owner is found, then the owner struct isn't populated, which gives ID=0, but ID's start at 1, so this shouldn't cause any potential security issues
	managerExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, newManagerUname, newManager)
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	if !ownerExists || !managerExists {
		s.Message = status.UserNotFound
		return http.StatusNotFound, nil
	}
	if !(isOwner(db, owner, club) && owner.Username != newManager.Username) {
		s.Message = http.StatusText(http.StatusForbidden)
		return http.StatusForbidden, nil
	}
	var err error
	switch op {
	case model.OpAdd:
		res := db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ?", newManager.ID, club.ID).First(model.NewUserClub())
		if res.Error != nil { // Record not existing, then add the new manager
			err = db.Model(newManager).Association(model.ManagesColumn).Append(club)
			if err != nil {
				s.Message = failureMessage
				return http.StatusInternalServerError, fmt.Errorf("unable to add new manager")
			}
		}
	case model.OpRemove:
		err = db.Model(newManager).Association(model.ManagesColumn).Delete(club)
		if err != nil {
			s.Message = failureMessage
			return http.StatusInternalServerError, fmt.Errorf("unable to remove manager")
		}
	}
	s.Message = successMessage
	s.Code = status.SuccessCode
	return http.StatusOK, nil
}

// LeaveClub lets the user step down a club manager (The user won't have any correlations previously managed club)
// Note: If the user is a club owner, they must appoint a new owner in replacement of them
func LeaveClub(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	club := model.NewClub()
	user := model.NewUser()
	clubID := getVar(r, model.ClubIDVar)
	claims := GetTokenClaims(ExtractToken(r))
	uname := fmt.Sprintf("%v", claims["sub"])
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	isOwner := isOwner(db, user, club)
	if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	if !userExists {
		s.Message = status.UserNotFound
		return http.StatusNotFound, nil
	}
	if isOwner {
		s.Message = status.LeaveClubFailure
		return http.StatusUnprocessableEntity, nil
	}
	if !isManager(db, user, club) {
		s.Message = http.StatusText(http.StatusForbidden)
		return http.StatusForbidden, nil
	}
	err := db.Model(user).Association(model.ManagesColumn).Delete(club)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("unable to leave club")
	}
	s.Code = status.SuccessCode
	s.Message = status.LeaveClubSuccess
	return http.StatusOK, nil
}

// PromoteOwner promotes a club manager to be the new owner while the current owner would step down and become a manager
// Note: There can only be 1 club owner but you can have many club managers
func PromoteOwner(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	club := model.NewClub()
	potentialNewOwner := model.NewUser()
	currentOwner := model.NewUser()
	clubID := getVar(r, model.ClubIDVar)
	potentialNewOwnerUname := getVar(r, model.UsernameVar)
	claims := GetTokenClaims(ExtractToken(r))
	currentOwnerUname := fmt.Sprintf("%v", claims["sub"])
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	newOwnerExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, potentialNewOwnerUname, potentialNewOwner)
	currentOwnerExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, currentOwnerUname, currentOwner)
	if !currentOwnerExists || !newOwnerExists {
		s.Message = status.UserNotFound
		return http.StatusNotFound, nil
	}
	if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	if !isOwner(db, currentOwner, club) {
		s.Message = status.ClubPromoteOwnerFailure
		return http.StatusUnprocessableEntity, nil
	}
	if !isManager(db, potentialNewOwner, club) {
		s.Message = status.ClubPromoteNeedManager
		return http.StatusUnprocessableEntity, nil
	}
	if potentialNewOwnerUname == currentOwnerUname {
		s.Message = status.ClubPromoteSelfFailure
		return http.StatusUnprocessableEntity, nil
	}
	// New owner assumes new owner position
	res := db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ? AND is_owner = ?", potentialNewOwner.ID, club.ID, false).Update(model.IsOwnerColumn, true)
	if res.Error != nil {
		return http.StatusInternalServerError, fmt.Errorf("unable to promote user to new owner")
	}
	// Old owner steps down and becomes manager
	res = db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ? AND is_owner = ?", currentOwner.ID, club.ID, true).Update(model.IsOwnerColumn, false)
	if res.Error != nil {
		return http.StatusInternalServerError, fmt.Errorf("unable to convert user from owner to manager")
	}
	s.Code = status.SuccessCode
	s.Message = status.ClubPromoteSuccess
	return http.StatusOK, nil
}
