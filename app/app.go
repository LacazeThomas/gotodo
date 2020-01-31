package app

import (
	"fmt"
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
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.Charset)

	db, err := gorm.Open(config.Dialect, dbURI)
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
	a.Post("/api/user/new", a.handleRequest(handler.CreateAccount))
	a.Post("/api/user/login", a.handleRequest(handler.Authenticate))

	// Routing for handling the projects
	a.Get("/projects", a.handleRequest(handler.GetAllProjects))
	a.Post("/projects", a.handleRequest(handler.CreateProject))
	a.Get("/projects/{title}", a.handleRequest(handler.GetProject))
	a.Put("/projects/{title}", a.handleRequest(handler.UpdateProject))
	a.Delete("/projects/{title}", a.handleRequest(handler.DeleteProject))
	a.Put("/projects/{title}/archive", a.handleRequest(handler.ArchiveProject))
	a.Delete("/projects/{title}/archive", a.handleRequest(handler.RestoreProject))

	// Routing for handling the tasks
	a.Get("/projects/{title}/tasks", a.handleRequest(handler.GetAllTasks))
	a.Post("/projects/{title}/tasks", a.handleRequest(handler.CreateTask))
	a.Get("/projects/{title}/tasks/{id:[0-9]+}", a.handleRequest(handler.GetTask))
	a.Put("/projects/{title}/tasks/{id:[0-9]+}", a.handleRequest(handler.UpdateTask))
	a.Delete("/projects/{title}/tasks/{id:[0-9]+}", a.handleRequest(handler.DeleteTask))
	a.Put("/projects/{title}/tasks/{id:[0-9]+}/complete", a.handleRequest(handler.CompleteTask))
	a.Delete("/projects/{title}/tasks/{id:[0-9]+}/complete", a.handleRequest(handler.UndoTask))
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