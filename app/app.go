package app

import (
	"../config"
	"./handler"
	"./logger"
	"./model"
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

type routeHandler func(w http.ResponseWriter, r *http.Request)
type hdlr func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

type App struct {
	db      *gorm.DB
	router  *mux.Router
	origin  handlers.CORSOption
	methods handlers.CORSOption
	headers handlers.CORSOption
}

func (app *App) Initialize(dbConfig *config.DBConfig, redisConfig *config.RedisConfig) {
	ctx := context.Background()
	dbFormat :=
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Name,
		)
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     redisConfig.Addr,
			Password: redisConfig.Password,
			DB:       redisConfig.DB,
		})
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Unable to connect to Redis\n", err)
	}

	db, err := gorm.Open("postgres", dbFormat)
	if err != nil {
		log.Fatal("Unable to connect to database\n", err)
	}
	app.db = db
	app.router = mux.NewRouter().StrictSlash(true)
	app.router.Use(logger.LoggingMiddleware)
	// Note: Set this as env var later
	app.origin = handlers.AllowedOrigins([]string{"http://localhost:3000"})
	app.methods = handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete})
	app.headers = handlers.AllowedHeaders([]string{"Content-Type"})

	app.setRoutes()
	log.Println("Connected to Redis")
	log.Println("Connected to Database")
	db.SingularTable(true)
	db.CreateTable(model.NewStudent(), model.NewEvent(), model.NewLog(), model.NewClub(), model.NewTag())

}

func (app *App) setRoutes() {
	// Signup Route
	app.Post("/signup", app.Handle(handler.SignUp, false)) // Done
	app.Get("/signup/usernames/{username}", app.Handle(handler.QueryUsername, false)) // Done
	app.Get("/signup/emails/{email}", app.Handle(handler.QueryEmail, false)) // Done

	// Login Routes
	app.Post("/login", app.Handle(handler.Login, false)) // Done (Might need to check for synchronous token)

	// Student Routes
	app.Get("/students/{username}", app.Handle(handler.GetStudent, true)) // Partially working (must only return partial information)
	app.Post("/students/{username}", app.Handle(handler.UpdateStudent, true)) // POST

	app.Post("/test/{username}", app.Handle(handler.Test, false))
	// Club routes
	app.Get("/clubs", app.Handle(handler.GetClubs, false))
	app.Get("/clubs/tags/{tag}", app.Handle(handler.GetClubsTag, false))
	app.Get("/clubs/{username}", app.Handle(handler.GetClub, false))
	app.Post("/clubs/{username}", app.Handle(handler.UpdateClub, true)) // POST
	app.Get("/clubs/events", app.Handle(handler.GetEvents, true))
	app.Get("/clubs/events/{username}", app.Handle(handler.GetEvent, true))
	app.Post("/clubs/events/{username}", app.Handle(handler.CreateEvent, true)) // POST
	app.Post("/clubs/events/{username}", app.Handle(handler.UpdateEvent, true)) // POST
	app.Delete("/clubs/events/{username}", app.Handle(handler.DeleteEvent, true))

	// Chat Routes

	// 404 Route
	app.router.NotFoundHandler = handler.NotFound()
}

func (app *App) Run(port string) {
	http.ListenAndServe(port, handlers.CORS(app.origin, app.methods, app.headers)(app.router))
}

func (app *App) Post(path string, f routeHandler) {
	app.router.HandleFunc(path, f).Methods(http.MethodPost)
}

func (app *App) Get(path string, f routeHandler) {
	app.router.HandleFunc(path, f).Methods(http.MethodGet)
}

func (app *App) Delete(path string, f routeHandler) {
	app.router.HandleFunc(path, f).Methods(http.MethodDelete)
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
