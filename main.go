package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

type Todo struct {
	Id            int64          `json:"id"`
	Name          string         `json:"name"`
	Description   string         `json:"description,omitempty"`
	CompletedDate mysql.NullTime `json:"completion_date,omitempty"`
	DueDate       mysql.NullTime `json:"due_date",omitempty`
	CreatedAt     mysql.NullTime `json:"created_at_date,omitempty"`
}

func main() {
	db, err := gorm.Open("mysql", "todo:some_pass@/todos?charset=utf8&parseTime=True")
	if err != nil {
		panic(err.Error())
	}

	if db.HasTable(&Todo{}) == false {
		fmt.Println("Creating todos table")

		err := db.CreateTable(&Todo{}).Error
		if err != nil {
			panic(err)
		}
	}

	err = db.AutoMigrate(&Todo{}).Error
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter().StrictSlash(false)
	todos := r.Path("/todos").Subrouter()
	todos.Methods("GET").Handler(ShowTodos(db))
	todos.Methods("POST").Handler(CreateTodo(db))

	todo := r.PathPrefix("/todos/{id}").Subrouter()
	todo.Methods("PUT").Path("/complete").Handler(TodoToggleComplete(db))
	todo.Methods("PUT").Handler(UpdateTodo(db))

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.ListenAndServe(":8080", r)
}

func ShowTodos(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		allTodos := []Todo{}
		err := db.Find(&allTodos).Error
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(allTodos)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(js)
	})
}

func CreateTodo(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var newTodo Todo
		err := decoder.Decode(&newTodo)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		err = db.Save(&newTodo).Error
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(newTodo)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(js)
	})
}

const missingIDErrorMessage string = "Missing or malformed id field"

func UpdateTodo(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idOfTodoToUpdate, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(rw, missingIDErrorMessage, http.StatusBadRequest)
			return
		}

		var existingTodo Todo
		err = db.First(&existingTodo, idOfTodoToUpdate).Error
		if err == gorm.ErrRecordNotFound {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var todoToUpdate Todo
		err = decoder.Decode(&todoToUpdate)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		err = db.Save(&todoToUpdate).Error
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(&todoToUpdate)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(js)
	})
}

func TodoToggleComplete(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idOfTodoToUpdate, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(rw, missingIDErrorMessage, http.StatusBadRequest)
			return
		}

		var todoToComplete Todo
		err = db.Find(&todoToComplete, idOfTodoToUpdate).Error
		if err == gorm.ErrRecordNotFound {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		if todoToComplete.CompletedDate.Valid == false {
			todoToComplete.CompletedDate.Time = time.Now()
		} else {
			todoToComplete.CompletedDate.Time = time.Time{}
		}

		err = db.Save(&todoToComplete).Error
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		js, err := json.Marshal(todoToComplete)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(js)
	})
}
