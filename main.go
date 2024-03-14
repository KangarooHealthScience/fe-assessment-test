package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

type Todo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Details   string    `json:"details"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

type Response struct {
	Status       string `json:"string"`
	ErrorMessage any    `json:"error_message,omitempty"`
	Data         any    `json:"data"`
}

var host = "0.0.0.0:3000"
var todos = make(map[string]Todo)
var mtx = new(sync.RWMutex)
var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

// @title KH FE Assessment Test
// @version 1.0
// @description This is the web service of simple TODO list app
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /api
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Post("/login", LoginUser)

			r.Route("/todo", func(r chi.Router) {
				r.Use(jwtauth.Verifier(tokenAuth))
				r.Use(jwtauth.Authenticator(tokenAuth))
				r.Get("/", GetTodoList)
				r.Post("/", AddTodoList)
				r.Put("/{todoID}", UpdateTodoList)
				r.Delete("/{todoID}", DeleteTodoList)
			})
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("welcome. navigate to http://%s/swagger/index.html to see the api docs", strings.ReplaceAll(host, "0.0.0.0", "localhost"))))
	})

	log.Println("starting server at", host)
	http.ListenAndServe(host, r)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		writeResponse(w, "error", nil, fmt.Errorf("unable to complete auth"))
		return
	}

	if !(username == "kangaroohealth" && password == "the magnificent chicken") {
		writeResponse(w, "error", nil, fmt.Errorf("invalid username and/or password"))
		return
	}

	claims := map[string]interface{}{
		"id":       uuid.New().String(),
		"username": username,
	}
	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		writeResponse(w, "error", nil, err)
		return
	}

	writeResponse(w, "ok", tokenString, nil)
}

func GetTodoList(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, "ok", todosAsSlice(), nil)
}

func AddTodoList(w http.ResponseWriter, r *http.Request) {
	todo := new(Todo)
	err := json.NewDecoder(r.Body).Decode(todo)
	if err != nil {
		writeResponse(w, "error", nil, err)
		return
	}

	if todo.Name == "" {
		writeResponse(w, "error", nil, fmt.Errorf("name cannot be empty"))
		return
	}

	todo.ID = uuid.New().String()
	todo.CreatedAt = time.Now()

	mtx.Lock()
	todos[todo.ID] = *todo
	mtx.Unlock()

	writeResponse(w, "ok", todosAsSlice(), nil)
}

func UpdateTodoList(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoID")
	if todoID == "" {
		writeResponse(w, "error", nil, fmt.Errorf("todoID cannot be empty"))
		return
	}

	patch := new(Todo)
	err := json.NewDecoder(r.Body).Decode(patch)
	if err != nil {
		writeResponse(w, "error", nil, err)
		return
	}
	if patch.Name == "" {
		writeResponse(w, "error", nil, fmt.Errorf("name cannot be empty"))
		return
	}

	mtx.RLock()
	v, ok := todos[todoID]
	mtx.RUnlock()

	if !ok {
		err = fmt.Errorf("todo with id %s does not exists", todoID)
		writeResponse(w, "error", nil, err)
		return
	}

	mtx.Lock()
	v.Name = patch.Name
	v.Details = patch.Details
	v.Done = patch.Done
	todos[todoID] = v
	mtx.Unlock()

	writeResponse(w, "ok", todosAsSlice(), nil)
}

func DeleteTodoList(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoID")
	if todoID == "" {
		writeResponse(w, "error", nil, fmt.Errorf("todoID cannot be empty"))
		return
	}

	mtx.RLock()
	_, ok := todos[todoID]
	mtx.RUnlock()

	if !ok {
		err := fmt.Errorf("todo with id %s does not exists", todoID)
		writeResponse(w, "error", nil, err)
		return
	}

	mtx.Lock()
	delete(todos, todoID)
	mtx.Unlock()

	writeResponse(w, "ok", todosAsSlice(), nil)
}

func writeResponse(w http.ResponseWriter, status string, data any, err error) {
	w.Header().Set("Content-type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(Response{
			Status:       status,
			Data:         data,
			ErrorMessage: err.Error(),
		})
	} else {
		json.NewEncoder(w).Encode(Response{
			Status: status,
			Data:   data,
		})
	}
}

func todosAsSlice() []Todo {
	mtx.RLock()
	res := make([]Todo, 0)
	for _, v := range todos {
		res = append(res, v)
	}
	mtx.RUnlock()

	sort.Slice(res, func(i, j int) bool {
		l := res[i].CreatedAt
		r := res[j].CreatedAt
		return l.UnixNano() < r.UnixNano()
	})

	return res
}
