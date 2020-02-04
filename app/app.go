package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"github.com/lacazethomas/goTodo/app/handler"
	"github.com/lacazethomas/goTodo/app/model"
	"github.com/lacazethomas/goTodo/config"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
	Token     string
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config config.DB) {

	db, err := gorm.Open("postgres", "host='192.168.1.35' port=5432 user=nextcloud dbname='gotodo-dev' password='YOUR_SECRET_PASSWORD' sslmode=disable")
	/* db, err := gorm.Open(config.Dialect, dbURI) */
	if err != nil {
		log.Fatal(err.Error())
	}

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()

	a.Router.Use(handler.JwtAuthentication)
	
	a.setRouters()
}

// setRouters sets the all required routers
func (a *App) setRouters() {
	
	// Routing for handling the login
	a.Post("/user/register", a.handleRequest(handler.CreateAccount))
	a.Post("/user/login", a.handleRequest(handler.Authenticate))

	// Routing for handling the projects
	a.Get("/projects/{status:[0-1]}", a.handleRequest(handler.GetAllProjects))
	a.Post("/project", a.handleRequest(handler.CreateProject))
	a.Get("/project/{id}", a.handleRequest(handler.GetProject))
	a.Put("/project/{id}", a.handleRequest(handler.UpdateProject))
	a.Delete("/project/{id}", a.handleRequest(handler.DeleteProject))
	a.Put("/project/{id}/archive", a.handleRequest(handler.ArchiveProject))
	a.Delete("/project/{id}/archive", a.handleRequest(handler.RestoreProject))
 
	// Routing for handling the tasks
	a.Get("/projects/{id}/tasks", a.handleRequest(handler.GetAllTasks))
	a.Post("/projects/{id}/tasks", a.handleRequest(handler.CreateTask))
	a.Get("/projects/{id}/tasks/{id:[0-9]+}", a.handleRequest(handler.GetTask))
	a.Put("/projects/{id}/tasks/{id:[0-9]+}", a.handleRequest(handler.UpdateTask))
	a.Delete("/projects/{id}/tasks/{id:[0-9]+}", a.handleRequest(handler.DeleteTask))
	a.Put("/projects/{id}/tasks/{id:[0-9]+}/complete", a.handleRequest(handler.CompleteTask))
	a.Delete("/projects/{id}/tasks/{id:[0-9]+}/complete", a.handleRequest(handler.UndoTask))
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}