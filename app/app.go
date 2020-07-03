package app

import (
	"fmt"
	"github.com/2-of-Clubs/2ofclubs-server/app/handler"
	"github.com/2-of-Clubs/2ofclubs-server/app/logger"
	"github.com/2-of-Clubs/2ofclubs-server/app/model"
	"github.com/2-of-Clubs/2ofclubs-server/config"
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
	//ctx := context.Background
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
	fmt.Println(redisClient)
	//pong, err := redisClient.Ping(ctx).Result()
	//fmt.Println(pong)

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
	log.Println("Connected to Database")
	db.SingularTable(true)
	db.CreateTable(model.NewStudent(), model.NewPerson(), model.NewEvent(), model.NewChat(), model.NewLog())

}

func (app *App) setRoutes() {

	// Login Routes
	app.Post("/login", app.Handle(handler.Login, true))
	// Student Routes
	app.Post("/signup", app.Handle(handler.CreateStudent, false))
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

func (app *App) Run(port string) {
	http.ListenAndServe(port, handlers.CORS(app.origin, app.methods, app.headers)(app.router))
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
