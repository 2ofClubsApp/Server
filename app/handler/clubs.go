package handler

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// In-Progress
func GetClubs(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	s := status.New()

	type Club struct {
		club_id int
	}
	var clubs []model.Club
	//clubs := []model.Club{}
	//s.Message = status.ClubsFound
	activeTags := flatten(filterTags(extractTags(db, r)))
	fmt.Println(activeTags)
	//db.Table(model.ClubTagTable).Where("tag_name IN ?", activeTags).Find(&clubs)
	db.Joins("JOIN club_tag ON club_tag.club_id=club.id").
		Joins("JOIN tag ON club_tag.tag_name=tag.name").
		Where("tag.name IN ?", activeTags).
		Distinct("club.name").
		Find(&clubs)
	//db.Raw("SELECT DISTINCT club.name FROM club, club_tag JOIN tag ON club_tag.tag_name=tag.name JOIN club_tag ON club_tag.club_id=club.id WHERE tag.name IN ?", activeTags).Scan(&clubs)
	//res := db.Raw("Select club.id From club NATURAL JOIN club_tag Where club_tag.tag_name IN ?", activeTags).Find(&clubs)
	//fmt.Println(res.RowsAffected)
	//fmt.Println(res.Error)
	//fmt.Println(clubs)
	for _, r := range clubs {
		fmt.Println(r.Name)
	}
	WriteData(GetJSON(s), http.StatusOK, w)
}

func UpdateClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	clubID := getVar(r, model.ClubIDVar)
	httpStatusCode := http.StatusOK
	club := model.NewClub()
	user := model.NewUser()
	s := status.New()
	claims := GetTokenClaims(r)
	uname := fmt.Sprintf("%v", claims["sub"])
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	if clubExists && userExists && isManager(db, user, club) {
		updatedClub := model.NewClub()
		extractBody(r, updatedClub)
		validate := validator.New()
		err := validate.Struct(updatedClub)
		if err == nil {
			db.Model(club).Select(model.BioColumn, model.SizeColumn).Updates(updatedClub)
			s.Code = status.SuccessCode
			s.Message = status.ClubUpdateSuccess
		} else {
			s.Message = status.ClubUpdateFailure
		}
	} else if !clubExists {
		s.Message = status.ClubNotFound
	} else {
		s.Code = status.FailureCode
		s.Message = http.StatusText(http.StatusForbidden)
		httpStatusCode = http.StatusForbidden
	}
	WriteData(GetJSON(s), httpStatusCode, w)
}

/*
Check if the email & username is available (RecordExists)
*/
func CreateClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	claims := GetTokenClaims(r)
	user := model.NewUser()
	uname := fmt.Sprintf("%v", claims["sub"])
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	club := model.NewClub()
	extractBody(r, club)
	validate := validator.New()
	err := validate.Struct(club)
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.NameColumn, club.Name, model.NewClub())
	emailExists := SingleRecordExists(db, model.ClubTable, model.EmailColumn, club.Email, model.NewClub())
	s := status.New()
	// Keeping userExists as a check even though the user should exist given the valid token because there's a chance that the user is deleted
	// In this case the user will still exist in the database but will be inaccessible.
	fmt.Println(userExists)
	fmt.Println(clubExists)
	fmt.Println(emailExists)
	if !emailExists && !clubExists && userExists && err == nil {
		db.Model(user).Association(model.ManagesColumn).Append(club)
		db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ? AND is_owner = ?", user.ID, club.ID, false).Update(model.IsOwnerColumn, true)
		s.Message = status.ClubCreationSuccess
		s.Code = status.SuccessCode
	} else {
		s.Message = status.ClubCreationFailure
	}
	WriteData(GetJSON(s), http.StatusOK, w)
}

func UpdateClubTags(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var httpStatus int
	s := status.New()
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	claims := GetTokenClaims(r)
	username := fmt.Sprintf("%v", claims["sub"])
	user := model.NewUser()
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, username, user)
	// Must check with both user and club existing in the event that a user gets deleted but you manage to get a hold of their access token
	fmt.Println("Checkpoint 1")
	if userExists && clubExists && isManager(db, user, club) {
		fmt.Println("Checkpoint 2")
		tags := filterTags(extractTags(db, r))
		db.Model(club).Association(model.SetsColumn).Replace(tags)
		s.Message = status.TagsUpdated
		s.Code = status.SuccessCode
		httpStatus = http.StatusOK
	} else if !clubExists {
		s.Message = status.ClubNotFound
		httpStatus = http.StatusOK
	} else {
		fmt.Println("Checkpoint 4")
		s.Message = http.StatusText(http.StatusForbidden)
		httpStatus = http.StatusForbidden
	}
	WriteData(GetJSON(s), httpStatus, w)
}

func GetClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	getClubInfo(db, w, r, model.AllClubInfo)
}

func getClubInfo(db *gorm.DB, w http.ResponseWriter, r *http.Request, infoType string) {
	var statusCode int
	var data string
	clubID := getVar(r, model.ClubIDVar)
	s := status.New()
	club := model.NewClub()
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	if !clubExists {
		s.Message = status.ClubNotFound
	} else {
		switch strings.ToLower(infoType) {
		case model.AllClubInfo:
			loadClubData(db, club)
			s.Data = club
		case model.AllClubEventsHost:
			clubEvents := make(map[string][]model.Event)
			db.Table(model.ClubTable).Preload(model.HostsColumn).Find(club)
			clubEvents[model.HostsColumn] = club.Hosts
			s.Data = clubEvents
		}
		s.Message = status.ClubFound
		s.Code = status.SuccessCode
	}
	statusCode = http.StatusOK
	data = GetJSON(s)
	WriteData(data, statusCode, w)
}

func GetClubEvents(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	getClubInfo(db, w, r, model.AllClubEventsHost)
}
func loadClubData(db *gorm.DB, club *model.Club) {
	db.Table(model.ClubTable).Preload(model.SetsColumn).Find(club)
	db.Table(model.ClubTable).Preload(model.HostsColumn).Find(club)
	club.Sets = filterTags(club.Sets)
	club.Hosts = club.Hosts
}

/*
Returns true if the user is an owner of the club, false otherwise
*/
func isOwner(db *gorm.DB, user *model.User, club *model.Club) bool {
	userClub := model.NewUserClub()
	db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ?", user.ID, club.ID).First(userClub)
	return userClub.IsOwner
}

func isManager(db *gorm.DB, user *model.User, club *model.Club) bool {
	userClub := model.NewUserClub()
	res := db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ?", user.ID, club.ID).First(userClub)
	return res.Error == nil
}

func RemoveManager(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	editManagers(db, w, r, model.OpRemove)
}

func AddManager(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	editManagers(db, w, r, model.OpAdd)
}

/*
Adding or removing managers and their associations to a particular club
*/
func editManagers(db *gorm.DB, w http.ResponseWriter, r *http.Request, op string) {
	// Default messages set to manager addition, otherwise manager removal
	var successMessage = status.ManagerAdditionSuccess
	var failureMessage = status.ManagerAdditionFailure
	if op == model.OpRemove {
		successMessage = status.ManagerRemoveSuccess
		failureMessage = status.ManagerRemoveFailure
	}
	s := status.New()
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
				} else {
					err = fmt.Errorf("unable to add manager")
				}
			case model.OpRemove:
				err = db.Model(newManager).Association(model.ManagesColumn).Delete(club)
			}
			if err != nil {
				s.Message = failureMessage
			} else {
				s.Message = successMessage
				s.Code = status.SuccessCode
			}
		} else {
			s.Message = failureMessage
		}
	} else {
		s.Message = failureMessage
	}
	WriteData(GetJSON(s), http.StatusOK, w)
}
