package handler

import (
	"../model"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"net/http"
)

const ErrGeneric = "an error occurred"

/*
	Common methods shared amongst the different models
*/

func NotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteData(http.StatusText(http.StatusNotFound), http.StatusNotFound, w)
	})
}
/*
Extract the Token Claims from the HTTP Request Header
*/
func GetTokenClaims(r *http.Request) jwt.MapClaims{
	t := r.Header["Token"][0]
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(t, &claims, kf)
	return claims
}

func Test(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	//c := model.NewClub()
	//t := model.NewTag()
	//t.Name = "Computer Science"

	//t1 := model.NewTag()
	//t1.Name = "Banana"

	//e := model.NewEvent()
	//e.DateTime = "now"
	//e.Description = "So fun much wow"
	//e.Fee = 0.99

	//e1 := model.NewEvent()
	//e1.DateTime = "tmrw"
	//e1.Description = "wow much fun so"
	//e1.Fee = 99.0
	//////////////////////
	u := model.NewUser()
	//u.Username="Hiimchrislim"
	//u.Password = "password"
	//u.Email = "hello@hiimchrislim.co"
	cc := model.NewClub()
	cc.Bio = "We are ACS!"
	cc.Size = 123456789
	cc.Email = "acs@utm.com"
	//u.Manages = []model.Club{*cc}
	///////////////////////////
	db.Table("user").Where("username = ?", "Hiimchrislim").First(&u)

	//db.Model(u).Association("Manages").Append([]model.Club{*cc})
	uc := model.UserClub{}
	db.Table("user_club").Where("user_id = ?", string(u.ID)).First(&uc)
	db.Table("user_club").Where("user_id = ? AND club_id = ?", u.ID).Update("is_owner", true)

	//cc.HelpNeeded = true
	//e := model.NewEvent()
	//var club [] model.Club
	//db.Preload(model.HostsColumn).Find(&club)
	//a := db.Model(c).Association("Hosts").Count()
	//fmt.Println(a)
	//db.Model(c).Table(model.ClubTable).Where("Username = ?", "Hacklab").Updates(*cc)
}

/*
Note: Need to add more authentication checks later (This is temporary)
*/
func IsValidJWT(w http.ResponseWriter, r *http.Request) bool {
	if token := r.Header["Token"]; token != nil {
		if t, err := jwt.Parse(token[0], kf); err == nil {
			if t.Valid {
				return true
			}
		}
	}
	WriteData(http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized, w)
	return false
}

func kf(token *jwt.Token) (interface{}, error) {
	// Verifying that the signing method is the same before continuing any further
	if _, accepted := token.Method.(*jwt.SigningMethodHMAC); !accepted {
		return nil, fmt.Errorf(ErrGeneric)
	}
	// Note: This must be changed to an env variable later
	return []byte("2ofClubs"), nil
}

/*
Returning true if the record already exists in the table, false otherwise.
*/
//You can put a check on Record Exists on the deleted column as long as it's null it'll exist then
// When soft deleted, the record won't exist anymore
func SingleRecordExists(db *gorm.DB, tableName string, column string, val string, t interface{}) bool {
	result := db.Table(tableName).Where(column+"= ?", val).First(t)
	return result.Error == nil
}

/*
Returning the JSON representation of a struct.
*/
func GetJSON(response interface{}) string {
	data, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

/*
Return response message and an HTTP Status Code upon receiving a request.
*/
func WriteData(data string, code int, w http.ResponseWriter) int {
	w.WriteHeader(code)
	n, err := fmt.Fprint(w, data)
	if err != nil {
		return -1
	}
	return n
}
