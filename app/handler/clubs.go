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

// GetClubs - In-Progress
func GetClubs(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	type Club struct {
		clubId int
	}
	var clubs []model.Club
	//clubs := []model.Club{}
	//s.Message = status.ClubsFound
	activeTags := flatten(filterTags(extractTags(db, r)))
	fmt.Println(activeTags)
	db.Table(model.ClubTagTable).Where("tag_name IN ?", activeTags).Find(&clubs)

	db.Joins("JOIN club_tag ON club_tag.club_id=club.id").
		Joins("JOIN tag ON club_tag.tag_name=tag.name").
		Where("tag.name IN ?", activeTags).
		Distinct("club.name").
		Find(&clubs)
	//db.Joins("club").Joins("club_tag").Joins("tags").Find(&clubs, "club.sets IN ?", activeTags)
	//db.Raw("SELECT DISTINCT c.name FROM club AS c NATURAL JOIN club_tag AS ct WHERE tag.name IN ?", activeTags).Scan(&clubs)
	//db.Raw(" SELECT * FROM club_tag WHERE tag_name IN ?", activeTags).Find(&clubs)
	//res := db.Raw("Select club.id From club NATURAL JOIN club_tag Where club_tag.tag_name IN ?", activeTags).Find(&clubs)
	//fmt.Println(res.RowsAffected)
	//fmt.Println(res.Error)
	fmt.Println(clubs)
	for _, r := range clubs {
		fmt.Println(r.Name)
	}
	return http.StatusForbidden, nil
}

// UpdateClub - Update club with new information
// See model.Club or docs for club attribute specifications
func UpdateClub(db *gorm.DB, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	user := model.NewUser()
	claims := GetTokenClaims(r)
	uname := fmt.Sprintf("%v", claims["sub"])
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	if clubExists && userExists && isManager(db, user, club) {
		updatedClub := model.NewClub()
		extractBody(r, updatedClub)
		validate := validator.New()
		err := validate.Struct(updatedClub)
		if err != nil {
			s.Message = status.ClubUpdateFailure
			return http.StatusUnprocessableEntity, nil
		}
		res := db.Model(club).Select(model.BioColumn, model.SizeColumn).Updates(updatedClub)
		if res.Error != nil {
			return http.StatusInternalServerError, fmt.Errorf("unable to update club")
		}
		s.Code = status.SuccessCode
		s.Message = status.ClubUpdateSuccess
		return http.StatusOK, nil
	} else if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	s.Code = status.FailureCode
	s.Message = http.StatusText(http.StatusForbidden)
	return http.StatusForbidden, nil
}

// CreateClub - Creating a club (You must have an active user account first)
// See model.Club or the docs for club information constraints
func CreateClub(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	claims := GetTokenClaims(r)
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
	if clubExists {
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
	s.Message = status.ClubNotFound
	return http.StatusNotFound, nil
}

// UploadClubPhoto - Uploading a club photo
// Club photo file size upload is 10 MB max
func UploadClubPhoto(db *gorm.DB,_ *redis.Client,  _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	maxMem := int64(10 << 20)
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	user := model.NewUser()
	claims := GetTokenClaims(r)
	username := fmt.Sprintf("%v", claims["sub"])
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, username, user)
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	if clubExists && userExists && isManager(db, user, club) {
		// Max 10MB upload file
		err := r.ParseMultipartForm(maxMem) // 2^20
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("photo is exceeding upload limit")
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
		fileName := fmt.Sprintf("./images/%v.png", club.ID)
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
		s.Code = status.SuccessCode
		s.Message = status.FileWriteSuccess
		return http.StatusOK, nil
	}
	s.Message = http.StatusText(http.StatusForbidden)
	return http.StatusForbidden, nil

}

// UpdateClubTags - Updating user club tags
// All old tags will be overrided with the new set of tags provided
func UpdateClubTags(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	claims := GetTokenClaims(r)
	username := fmt.Sprintf("%v", claims["sub"])
	user := model.NewUser()
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, username, user)
	// Must check with both user and club existing in the event that a user gets deleted but you manage to get a hold of their access token
	if userExists && clubExists && isManager(db, user, club) {
		tags := filterTags(extractTags(db, r))
		err := db.Model(club).Association(model.SetsColumn).Replace(tags)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("unable to get club tags")
		}
		s.Message = status.TagsUpdated
		s.Code = status.SuccessCode
		return http.StatusOK, nil
	}
	if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	s.Message = http.StatusText(http.StatusForbidden)
	return http.StatusForbidden, nil
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
	case model.AllClubEventsHost:
		clubEvents := make(map[string][]model.Event)
		res := db.Table(model.ClubTable).Preload(model.HostsColumn).Find(club)
		if res.Error != nil {
			return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
		}
		clubEvents[strings.ToLower(model.HostsColumn)] = club.Hosts
		s.Data = clubEvents
	}
	s.Message = status.ClubFound
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
	claims := GetTokenClaims(r)
	clubOwnerUsername := fmt.Sprintf("%v", claims["sub"])
	newManagerUname := getVar(r, model.UsernameVar)
	clubID := getVar(r, model.ClubIDVar)
	owner := model.NewUser()
	newManager := model.NewUser()
	club := model.NewClub()
	// Added user must exist
	ownerExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, clubOwnerUsername, owner)
	// If owner is found, then the owner struct isn't populated, which gives ID=0, but ID's start at 1, so this shouldn't cause any potential security issues
	managerExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, newManagerUname, newManager)
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	if ownerExists && managerExists && clubExists {
		if isOwner(db, owner, club) && owner.Username != newManager.Username {
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
		s.Message = http.StatusText(http.StatusForbidden)
		return http.StatusForbidden, nil
	}
	s.Message = fmt.Sprintf("%s & %s", status.UserNotFound, status.ClubNotFound)
	return http.StatusNotFound, nil
}
