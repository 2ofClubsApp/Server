package handler

import (
	"../model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	ErrGeneric  = "an error occurred"
	ErrTokenGen = "token generation error"
	ErrSignUp   = "Unable to Sign Up Student"
	ErrLogin    = "Username or Password is Incorrect"
)

/*
	Common methods shared amongst the different models
*/
func SignUp(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Check if content type is application/json?
	u, isValid := VerifyUserInfo(r)
	status := model.NewStatus()
	status.Message = ErrSignUp
	if isValid {
		s := model.NewStudent()
		username := !RecordExists(db, model.StudentTable, model.UsernameColumn, u.Username, s)
		email := !RecordExists(db, model.StudentTable, model.EmailColumn, u.Email, s)
		if username && email {
			CreateStudent(db, w, u, s)
		} else {
			WriteData(GetJSON(status), http.StatusOK, w)
		}
	} else {
		WriteData(GetJSON(status), http.StatusOK, w)
	}
}

func NotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteData(http.StatusText(http.StatusNotFound), http.StatusNotFound, w)
	})
}

func Test(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	//c := model.NewClub()
	t := model.NewTag()
	t.Name = "Computer Science"

	t1 := model.NewTag()
	t1.Name = "Banana"

	e := model.NewEvent()
	e.DateTime = "now"
	e.Description = "So fun much wow"
	e.Fee = 0.99

	e1 := model.NewEvent()
	e1.DateTime = "tmrw"
	e1.Description = "wow much fun so"
	e1.Fee = 99.0

	cc := model.NewClub()
	cc.Username = "Banana"
	cc.Bio = "We are ACS!"
	cc.Password = "Hackhackhack"
	cc.Size = 123456789
	cc.Email = "acs@utm.com"
	cc.Hosts = []model.Event{*e, *e1}
	cc.Tags = []model.Tag{*t, *t1}
	db.Create(cc)

	//cc.HelpNeeded = true
	//e := model.NewEvent()
	//var club [] model.Club
	//db.Preload(model.HostsColumn).Find(&club)
	//a := db.Model(c).Association("Hosts").Count()
	//fmt.Println(a)
	//db.Model(c).Table(model.ClubTable).Where("Username = ?", "Hacklab").Updates(*cc)
}

func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	u, isValid := VerifyUserInfo(r)
	if isValid {
		hash, isFound := getPasswordHash(db, u.Username)
		err := errors.New("unable to find password")
		fmt.Println(isFound)
		if isFound {
			err = bcrypt.CompareHashAndPassword(hash, []byte(u.Password))
		}
		if err != nil {
			s := model.NewStatus()
			s.Message = ErrLogin
			WriteData(GetJSON(s), http.StatusOK, w)
		} else {
			if tp, err := GetTokenPair(u.Username, 5, 60*24); err == nil {
				c := GenerateCookie(model.RefreshToken, tp.RefreshToken)
				http.SetCookie(w, c)
				WriteData(tp.AccessToken, http.StatusOK, w)
			}
		}
	}
}

/*
	Gets password hash for both clubs and students provided the username.
*/
func getPasswordHash(db *gorm.DB, userName string) ([]byte, bool) {
	type p struct {
		Password string
	}
	pass := &p{}
	notFoundStudent := db.Table(model.StudentTable).Where("Username = ?", userName).Find(pass)
	if notFoundStudent != nil {
		return []byte(pass.Password), true
	}
	notFoundClub := db.Table(model.ClubTable).Where("Username = ?", userName).Find(pass)
	if notFoundClub != nil {
		return []byte(pass.Password), true
	}
	return []byte(""), false
}

/*
Validating the user request to ensure that they can only access/modify their own data.
 */
func ValidateUserReq(username string, r *http.Request) bool {
	t := r.Header["Token"][0]
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(t, &claims, kf)
	sub := fmt.Sprintf("%v", claims["sub"])
	fmt.Println(username)
	fmt.Println(sub)
	return sub == username
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
Generating http cookie where the refresh token will be embedded.
 */
func GenerateCookie(name string, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   true,
	}
}

func GenerateJWT(subject string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	sysTime := time.Now()
	claims["iat"] = sysTime
	claims["exp"] = sysTime.Add(time.Minute * duration).Unix()
	claims["sub"] = subject // Subject usually as a number (unique value?)
	// Note: This must be changed to an env variable later
	tokenString, err := token.SignedString([]byte("2ofClubs"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetTokenPair(subject string, accessDuration time.Duration, refreshDuration time.Duration) (*model.TokenInfo, error) {
	if accessToken, atErr := GenerateJWT(subject, accessDuration); atErr == nil {
		if refreshToken, rtErr := GenerateJWT(subject, refreshDuration); rtErr == nil {
			token := model.NewTokenInfo()
			token.AccessToken = accessToken
			token.RefreshToken = refreshToken
			return token, nil
		}
	}
	return nil, fmt.Errorf(ErrTokenGen)
}

/*
Returning (hash, true) on Hash success otherwise, ("", false) on error.
 */
func Hash(info string) (string, bool) {
	// Change cost to 10+ (try to find a way to scale it with hardware?)
	saltedHashPass, err := bcrypt.GenerateFromPassword([]byte(info), bcrypt.DefaultCost)
	if err != nil {
		return "", false
	}
	return string(saltedHashPass), true
}

/*
Returning true if the record already exists in the table, false otherwise.
 */
func RecordExists(db *gorm.DB, tableName string, column string, val string, t interface{}) bool {
	if db.Table(tableName).Where(column+"= ?", val).First(t) != nil {
		return true
	}
	return false
}

/*
Extracting JSON payload and returning (model, true) if valid, otherwise (model, false).
*/
func VerifyUserInfo(r *http.Request) (*model.User, bool) {
	decoder := json.NewDecoder(r.Body)
	u := model.NewUser()
	decoder.Decode(u)
	validate := validator.New()
	validate.RegisterValidation("alpha", ValidateUsername)
	err := validate.Struct(u)
	//fmt.Println(err)
	if err != nil {
		return u, false
	}
	u.Username = strings.ToLower(u.Username)
	u.Email = strings.ToLower(u.Email)
	return u, true
}

/*
Validate username against Regex pattern of being alphanumeric.
 */
func ValidateUsername(fl validator.FieldLevel) bool {
	matched, _ := regexp.Match("^[a-zA-Z0-9]+$", []byte(fl.Field().String()))
	//fmt.Printf("Valid Username: %v\n", matched)
	return matched
}

/*
Returning the representation of a struct formatted in JSON.
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
