package app

import (
	"../config"
	"../models"
	"./handlers"
	"./logger"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

type routeHandler func(w http.ResponseWriter, r *http.Request)
type handler func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

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
	app.router = mux.NewRouter()
	app.router.Use(logger.LoggingMiddleware)
	app.setRoutes()
	fmt.Println("Database Online")
	db.SingularTable(true)
	db.CreateTable(models.NewStudent(), models.NewPerson(), models.NewEvent(), models.NewChat(), models.NewLog())
}

func (app *App) setRoutes() {
	// Student Routes
	app.Post("/students", app.Handle(handlers.CreateStudent))
	app.Get("/students/{username}", app.Handle(handlers.GetStudent))
	app.Put("/students/{username}", app.Handle(handlers.UpdateStudent))

	// Club routes
	app.Get("/clubs", app.Handle(handlers.GetClubs))
	app.Post("/clubs", app.Handle(handlers.CreateClub))
	app.Get("/clubs/tags/{tag}", app.Handle(handlers.GetClubsTag))
	app.Get("/clubs/{username}", app.Handle(handlers.GetClub))
	app.Put("/clubs/{username}", app.Handle(handlers.UpdateClub))
	app.Get("/clubs/events", app.Handle(handlers.GetEvents))
	app.Get("/clubs/events/{username}", app.Handle(handlers.GetEvent))
	app.Post("/clubs/events/{username}", app.Handle(handlers.CreateEvent))
	app.Put("/clubs/events/{username}", app.Handle(handlers.UpdateEvent))
	app.Delete("/clubs/events/{username}", app.Handle(handlers.DeleteEvent))

	// Chat Routes

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

func (app *App) Handle(h handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h(app.db, w, r)
	}
}
