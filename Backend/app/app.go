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
	app.router = mux.NewRouter()
	app.router.Use(logger.LoggingMiddleware)
	app.setRoutes()
	fmt.Println("Database Online")
	db.SingularTable(true)
	db.CreateTable(model.NewStudent(), model.NewPerson(), model.NewEvent(), model.NewChat(), model.NewLog())
}

func (app *App) setRoutes() {
	// Student Routes
	app.Post("/students", app.Handle(handler.CreateStudent))
	app.Get("/students/{username}", app.Handle(handler.GetStudent))
	app.Put("/students/{username}", app.Handle(handler.UpdateStudent))

	// Club routes
	app.Get("/clubs", app.Handle(handler.GetClubs))
	app.Post("/clubs", app.Handle(handler.CreateClub))
	app.Get("/clubs/tags/{tag}", app.Handle(handler.GetClubsTag))
	app.Get("/clubs/{username}", app.Handle(handler.GetClub))
	app.Put("/clubs/{username}", app.Handle(handler.UpdateClub))
	app.Get("/clubs/events", app.Handle(handler.GetEvents))
	app.Get("/clubs/events/{username}", app.Handle(handler.GetEvent))
	app.Post("/clubs/events/{username}", app.Handle(handler.CreateEvent))
	app.Put("/clubs/events/{username}", app.Handle(handler.UpdateEvent))
	app.Delete("/clubs/events/{username}", app.Handle(handler.DeleteEvent))

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

func (app *App) Handle(h hdlr) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h(app.db, w, r)
	}
}
