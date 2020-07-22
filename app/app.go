package app

import (
	"../config"
	"./handler"
	"./logger"
	"./model"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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

func (app *App) Initialize(dbConfig *config.DBConfig, redisConfig *config.RedisConfig, adminConfig *model.User) {
	//ctx := context.Background()
	dbFormat :=
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Name,
		)
	//redisClient := redis.NewClient(
	//	&redis.Options{
	//		Addr:     redisConfig.Addr,
	//		Password: redisConfig.Password,
	//		DB:       redisConfig.DB,
	//	})
	//_, err := redisClient.Ping(ctx).Result()
	//if err != nil {
	//	log.Fatal("Unable to connect to Redis\n", err)
	//}

	db, err := gorm.Open(postgres.Open(dbFormat), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		//DisableForeignKeyConstraintWhenMigrating: true,
	})
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
	db.Migrator().CreateTable(model.NewEvent(), model.NewTag(), model.NewUserClub())
	db.SetupJoinTable(&model.User{}, "Manages", &model.UserClub{})
	db.AutoMigrate(model.NewUser(), model.NewClub())

}

func (app *App) setRoutes() {
	// Signup Route
	app.Post("/signup", app.Handle(handler.SignUp, false)) // Done
	//app.Get("/signup/usernames/{username}", app.Handle(handler.QueryUsername, false)) // Integrated into /signup
	//app.Get("/signup/emails/{email}", app.Handle(handler.QueryEmail, false))          // Integrated into /signup

	// Login Routes
	app.Post("/login", app.Handle(handler.Login, false)) // Done (Need to check for synchronous token (CSRF prevention))

	// User Routes
	app.Get("/users/{username}", app.Handle(handler.GetUser, true))              // Done
	app.Post("/users/{username}/tags", app.Handle(handler.UpdateUserTags, true)) // Done

	// Test Routes
	app.Post("/test/{username}", app.Handle(handler.Test, false)) // Ignore

	// Potential code merger on /clubs/{name} and /users/{username}

	// Tag Routes
	app.Get("/tags", app.Handle(handler.GetTags, false))               // Done
	app.Post("/tags/{tag}", app.Handle(handler.CreateTag, true))       // Done
	app.Delete("/tags/{tag}", app.Handle(handler.DeleteTag, true))     // Done
	app.Post("/upload/tags", app.Handle(handler.UploadTagsList, true)) // Done

	// Club routes
	app.Post("/clubs", app.Handle(handler.CreateClub, true))     // Done
	app.Get("/clubs/{name}", app.Handle(handler.GetClub, false)) // Done

	app.Delete("/clubs/{name}", app.Handle(handler.DeleteClub, true)) // Partially Done (The owner can delete the club and all associations will be removed?)
	app.Post("/clubs/{clubname}/manages/{username}", app.Handle(handler.AddManager, true)) // Done (Adding managers/maintainers to club)
	app.Delete("/clubs/{clubname}/manages/{username}", app.Handle(handler.RemoveManager, true)) // Partially done (Removing managers/maintainers) (If the current owner wants to leave, then they must appoint a new person)
	app.Post("/clubs/{clubname}/tags", app.Handle(handler.UpdateClubTags, true)) // (Adding tags for clubs)
	app.Get("/clubs", app.Handle(handler.GetClubs, false)) // In-Progress

	app.Get("/clubs/tags/{tag}", app.Handle(handler.GetClubsTag, false))
	app.Post("/clubs/{username}", app.Handle(handler.UpdateClub, true)) // POST
	app.Get("/clubs/events", app.Handle(handler.GetEvents, true))
	app.Get("/clubs/events/{username}", app.Handle(handler.GetEvent, true))
	app.Post("/clubs/events/{username}", app.Handle(handler.CreateEvent, true)) // POST
	app.Post("/clubs/events/{username}", app.Handle(handler.UpdateEvent, true)) // POST
	app.Delete("/clubs/events/{username}", app.Handle(handler.DeleteEvent, true))

	// Admin Route
	// Approve usernames
	// 404 Route
	app.router.NotFoundHandler = handler.NotFound() // Done
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
