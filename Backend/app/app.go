package app

import (
	"../config"
	"./handler"
	"./logger"
	"./model"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

type routeHandler func(w http.ResponseWriter, r *http.Request)
type hdlr func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

type App struct {
	db     *gorm.DB
	router *mux.Router
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
	app.db = db
	app.router = mux.NewRouter().StrictSlash(true)
	app.router.Use(logger.LoggingMiddleware)
	//app.router.Use(handler.ValidateJWT)
	app.setRoutes()
	log.Println("Connected to Database")
	db.SingularTable(true)
	db.CreateTable(model.NewStudent(), model.NewPerson(), model.NewEvent(), model.NewChat(), model.NewLog())
}

func (app *App) setRoutes() {
	// Student Routes
	app.Post("/students", app.Handle(handler.CreateStudent, false))
	app.Get("/students/{username}", app.Handle(handler.GetStudent, true))
	app.Put("/students/{username}", app.Handle(handler.UpdateStudent, true))

	// Club routes
	app.Get("/clubs", app.Handle(handler.GetClubs, true))
	app.Post("/clubs", app.Handle(handler.CreateClub, false))
	app.Get("/clubs/tags/{tag}", app.Handle(handler.GetClubsTag, true))
	app.Get("/clubs/{username}", app.Handle(handler.GetClub, true))
	app.Put("/clubs/{username}", app.Handle(handler.UpdateClub, true))
	app.Get("/clubs/events", app.Handle(handler.GetEvents, true))
	app.Get("/clubs/events/{username}", app.Handle(handler.GetEvent, true))
	app.Post("/clubs/events/{username}", app.Handle(handler.CreateEvent, true))
	app.Put("/clubs/events/{username}", app.Handle(handler.UpdateEvent, true))
	app.Delete("/clubs/events/{username}", app.Handle(handler.DeleteEvent, true))

	// Chat Routes

	// 404 Route
	app.router.NotFoundHandler = handler.NotFound()
}

func (app *App) Post(path string, f routeHandler) {
	app.router.HandleFunc(path, f).Methods(http.MethodPost)
}

func (app *App) Get(path string, f routeHandler) {
	app.router.HandleFunc(path, f).Methods(http.MethodGet)
}

func (app *App) Put(path string, f routeHandler) {
	app.router.HandleFunc(path, f).Methods(http.MethodPut)
}

func (app *App) Delete(path string, f routeHandler) {
	app.router.HandleFunc(path, f).Methods(http.MethodDelete)
}

func (app *App) Run(port string) {
	http.ListenAndServe(port, app.router)
}

func (app *App) Handle(h hdlr, verifyRequest bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Must verify for sensitive information
		if verifyRequest {
			if isValid := handler.IsValidJWT(w, r); isValid {
				h(app.db, w, r)
			}
		} else {
			h(app.db, w, r)
		}
	}
}
