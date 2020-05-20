package app

import (
	"../config"
	"../models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"time"
)

type App struct {
	db *gorm.DB
}

type test struct {
	//gorm.Model
	Name string
}

func (app *App) Initialize(config *config.DBConfig) {
	dbFormat :=
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Host,
			config.Port,
			config.User,
			config.Password,
			config.Name,
		)
	db, err := gorm.Open("postgres", dbFormat)
	if err != nil {
		log.Fatal("Unable to connect to database\n", err)
	}
	//defer db.Close()
	app.db = db
	fmt.Println("Connected")
	//fmt.Println(*db.Find(&models.Person{}))
	//db.SingularTable(true)
	//db.DropTableIfExists(&models.Person{})
	//db.DropTableIfExists(&test{})
	//db.CreateTable(&test{})
	//t := test{
	//	Name: "Sam",
	//}
	//db.Create(&t)
	//t.Name = "Bob"
	//db.Save(&t)
	//
	//var tt test
	//db.Debug().First(&tt)
	//fmt.Println(&tt)
	//db.Delete(&t)
	//tt = test{}
	//db.Debug().First(&tt)
	//fmt.Println(&tt)
	//db.CreateTable(test{})
	//db.Create(&t)
	type Author struct {
		gorm.Model
		FirstName string
		LastName  string
	}
	type Book struct {
		//gorm.Model
		Name        string
		PublishDate time.Time
		OwnerID     uint     `sql:"index"`
		Authors     []Author `gorm:"many2many:books_authors"`
	}
	type Owner struct {
		gorm.Model
		FirstName string
		LastName  string
		Books     []Book
	}
	p := models.Person{1, "Chris", "hello@hiimchrislim.co", "password"}
	t := models.Tags{[]string{"CS", "Math"}}
	b := models.Student{p, false, t}
	fmt.Printf("Type %T", b)
	db.DropTableIfExists(&Owner{}, &Book{}, &Author{})
	db.CreateTable(&Owner{}, &Book{}, &Author{}, &models.Stud	ent{})
}

/*
Test
 */
//tag := models.Tags{[]string{"CS", "Math", "Board Games"}}
////u := models.User{"hiimchrislim", "hello@hiimchrislim.co", "password"}
//u := models.User{}
//u.SetEmail("hello@hiimchrislim.co")
//u.SetUsername("hiimchrislim")
//u.SetPassword("password")
//m := models.Member{u, true, tag}
//fmt.Printf("S is %T\n", m)
//fmt.Printf("Email %s\n", m.GetEmail())
//fmt.Printf("Password %s\n", m.GetPassword())
//fmt.Printf("Username %s\n", m.GetUsername())
//hasTag, i := m.Tags.HasTag("Board Games")
//fmt.Printf("Tags %v\n @ index %d\n", hasTag, i)
