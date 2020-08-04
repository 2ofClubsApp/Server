package app

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/handler"
	"github.com/2-of-clubs/2ofclubs-server/app/logger"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/config"
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
	app.router = mux.NewRouter()
	//StrictSlash(true)
	app.router.Use(logger.LoggingMiddleware)
	// Note: Set this as env var later
	app.origin = handlers.AllowedOrigins([]string{"http://localhost:3000"})
	app.methods = handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions})
	app.headers = handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	app.setRoutes()
	log.Println("Connected to Redis")
	log.Println("Connected to Database")
	db.Migrator().CreateTable(model.NewEvent(), model.NewTag(), model.NewUserClub(), model.NewEvent())
	db.AutoMigrate(model.NewUser(), model.NewClub())
	db.SetupJoinTable(&model.User{}, "Manages", &model.UserClub{})

	// GORM already ensures the uniqueness of the username and email, thus we don't need to check if the admin already exists or not
	db.Create(adminConfig)
}

func (app *App) setRoutes() {
	// Signup Route
	app.Post("/signup", app.Handle(handler.SignUp, false)) // Done
	//app.Get("/signup/usernames/{username}", app.Handle(handler.QueryUsername, false)) // Integrated into /signup
	//app.Get("/signup/emails/{email}", app.Handle(handler.QueryEmail, false))          // Integrated into /signup

	// Login Routes
	app.Post("/login", app.Handle(handler.Login, false)) // Done (Need to check for synchronous token (CSRF prevention))

	// User Routes
	app.Get("/users/{username}", app.Handle(handler.GetUser, true))                     // Done
	app.Post("/users/{username}/tags", app.Handle(handler.UpdateUserTags, true))        // Done
	app.Get("/users/{username}/manages", app.Handle(handler.GetUserClubsManage, true))  // Done
	app.Get("/users/{username}/attends", app.Handle(handler.GetUserEventsAttend, true)) // Done
	app.Post("/events/{eid:[0-9]+}/attend", app.Handle(handler.AddUserAttendsEvent, true)) // Done
	app.Delete("/events/{eid:[0-9]+}/attend", app.Handle(handler.RemoveUserAttendsEvent, true)) // Done
	// Test Routes
	app.Post("/test/{username}", app.Handle(handler.Test, false)) // Ignore

	// Potential code merger on /clubs/{name} and /users/{username}

	// Tag Routes
	app.Get("/tags", app.Handle(handler.GetTags, false))                 // Done
	app.Post("/tags", app.Handle(handler.CreateTag, true))               // Done
	app.Post("/upload/tags", app.Handle(handler.UploadTagsList, true))   // Done
	app.Post("/tags/{tagName}/toggle", app.Handle(handler.ToggleTag, true)) // Done

	// Club routes
	app.Post("/clubs", app.Handle(handler.CreateClub, true))           // Done
	app.Post("/clubs/{cid:[0-9]+}", app.Handle(handler.UpdateClub, true)) // In-Progress
	app.Get("/clubs/{cid:[0-9]+}", app.Handle(handler.GetClub, false)) // Done

	//app.Delete("/clubs/{name}", app.Handle(handler.DeleteClub, true)) // Partially Done (The owner can delete the club and all associations will be removed?) (Clubs can't be deleted, only deactivated)
	app.Post("/clubs/{cid:[0-9]+}/manages/{username}", app.Handle(handler.AddManager, true))      // Done (Adding managers/maintainers to club)
	app.Delete("/clubs/{cid:[0-9]+}/manages/{username}", app.Handle(handler.RemoveManager, true)) // Partially done (Removing managers/maintainers) (If the current owner wants to leave, then they must appoint a new person)
	app.Post("/clubs/{cid:[0-9]+}/tags", app.Handle(handler.UpdateClubTags, true))                // Done
	app.Get("/clubs", app.Handle(handler.GetClubs, false))                                        // In-Progress
	//app.Get("/clubs/tags/{tag}", app.Handle(handler.GetClubsTag, false)) // Integrated into /clubs




	app.Get("/events", app.Handle(handler.GetAllEvents, false)) // Done
	app.Get("/events/{eid:[0-9]+}", app.Handle(handler.GetEvent, false)) // Done
	app.Get("/clubs/{cid:[0-9]+}/events", app.Handle(handler.GetClubEvents, false)) // Done

	app.Post("/clubs/{cid:[0-9]+}/events", app.Handle(handler.CreateClubEvent, true)) // Done

	app.Post("/clubs/{cid:[0-9]+}/events/{eid:[0-9]+}", app.Handle(handler.UpdateClubEvent, true)) // In-Progress
	app.Delete("/clubs/{cid:[0-9]+}/events/{eid:[0-9]+}", app.Handle(handler.DeleteClubEvent, true)) // Done

	// Admin Route
	app.Post("/users/{username}/toggle", app.Handle(handler.ToggleUser, true)) // Done
	//app.Post("/clubs/{cid}/toggle", app.Handle()) // In-Progress

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
