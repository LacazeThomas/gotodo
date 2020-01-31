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
	Base   *mux.Router
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
	a.Base = mux.NewRouter()

	a.Router = a.Base.PathPrefix("/api/projects").Subrouter().StrictSlash(true)
	a.setRouters()
}

// setRouters sets the all required routers
func (a *App) setRouters() {
	// Routing for handling the login
	a.PostLogin("/api/user/new", a.handleRequest(handler.CreateAccount))
	a.PostLogin("/api/user/login", a.handleRequest(handler.Authenticate))

	// Routing for handling the projects
	a.Get("", a.handleRequest(handler.GetAllProjects))
	a.Post("", a.handleRequest(handler.CreateProject))
	a.Get("/{title}", a.handleRequest(handler.GetProject))
	a.Put("/{title}", a.handleRequest(handler.UpdateProject))
	a.Delete("/{title}", a.handleRequest(handler.DeleteProject))
	a.Put("/{title}/archive", a.handleRequest(handler.ArchiveProject))
	a.Delete("/{title}/archive", a.handleRequest(handler.RestoreProject))

	// Routing for handling the tasks
	a.Get("/{title}/tasks", a.handleRequest(handler.GetAllTasks))
	a.Post("/{title}/tasks", a.handleRequest(handler.CreateTask))
	a.Get("/{title}/tasks/{id:[0-9]+}", a.handleRequest(handler.GetTask))
	a.Put("/{title}/tasks/{id:[0-9]+}", a.handleRequest(handler.UpdateTask))
	a.Delete("/{title}/tasks/{id:[0-9]+}", a.handleRequest(handler.DeleteTask))
	a.Put("/{title}/tasks/{id:[0-9]+}/complete", a.handleRequest(handler.CompleteTask))
	a.Delete("/{title}/tasks/{id:[0-9]+}/complete", a.handleRequest(handler.UndoTask))
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// PostLogin wraps the router for POST method
func (a *App) PostLogin(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Base.HandleFunc(path, f).Methods("POST")
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
	log.Fatal(http.ListenAndServe(host, a.Base))
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}

/*
func CreateToken() (string, error) {
	signingKey := []byte("keymaker")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": "hello",
		"role": "redpill",
	})
	tokenString, err := token.SignedString(signingKey)

	return tokenString, err
}

func ValidateToken(tokenString string) (jwt.Claims, error) {
	signingKey := []byte("keymaker")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		if token, ok := r.Header["Authorization"]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Authorization Header Missing")
		} else {
			if len(token) < 1 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Authorization Token Not Available")
			} else {
				_, err := ValidateToken(token[0])
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Authorization Token Not Valid")
				} else {
					// Call the next handler, which can be another middleware in the chain, or the final handler.
					next.ServeHTTP(w, r)
				}
			}
		}
	})
}

func LoginHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if token, err := CreateToken(); err != nil {
			fmt.Fprintf(w, "Error in Token Creation")
		} else {
			testmap := make(map[string]string)
			testmap["token"] = token
			b, err := json.Marshal(testmap)
			if err != nil {
				fmt.Fprintf(w, "Error in JSON Marshalling")
			}

			fmt.Fprintf(w, string(b))
		}
	}

	return http.HandlerFunc(fn)
}
*/
