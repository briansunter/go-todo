package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type Task struct {
	Id            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	CompletedDate time.Time `json:"completion_date"`
	DueDate       time.Time `json:"due_date"`
	CreatedAt     time.Time `json:"created_at_date"`
}

func main() {
	db, err := gorm.Open("mysql", "root:root@/todo?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	db.AutoMigrate(&Task{})

	r := mux.NewRouter().StrictSlash(false)
	tasks := r.Path("/tasks").Subrouter()
	tasks.Methods("GET").Handler(ShowTasks(&db))
	tasks.Methods("POST").Handler(CreateTask(&db))

	task := r.PathPrefix("/tasks/{id}").Subrouter()
	task.Methods("PUT").Path("/complete").Handler(TaskToggleComplete(&db))
	task.Methods("PUT").Handler(UpdateTask(&db))

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.ListenAndServe(":8080", r)
}

func ShowTasks(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		allTasks := []Task{}
		err := db.Find(&allTasks).Error
		if err != nil && err != gorm.RecordNotFound {
			panic(err)
		}

		js, err := json.Marshal(allTasks)
		if err != nil {
			panic(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(js)
	})
}

func CreateTask(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var newTask Task

		err := decoder.Decode(&newTask)

		if err != nil {
			http.Error(rw, err.Error(), 403)
			return
		}

		err = db.Save(&newTask).Error
		if err != nil {
			http.Error(rw, err.Error(), 403)
			return
		}

		js, err := json.Marshal(newTask)

		if err != nil {
			panic(err)
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(js)

	})
}

func UpdateTask(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		idOfTaskToUpdate := vars["id"]
		var existingTask Task
		err := db.First(&existingTask, idOfTaskToUpdate).Error

		if err == gorm.RecordNotFound {
			http.Error(rw, err.Error(), 404)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var taskToUpdate Task

		err = decoder.Decode(&taskToUpdate)

		if err != nil {
			http.Error(rw, err.Error(), 403)
			return
		}

		err = db.Save(&taskToUpdate).Error

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		js, err := json.Marshal(&taskToUpdate)

		if err != nil {
			panic(err)

		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(js)

	})

}

func TaskToggleComplete(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		idOfTaskToUpdate := vars["id"]

		var taskToComplete Task
		err := db.Find(&taskToComplete, idOfTaskToUpdate).Error
		if err != nil {

			if err == gorm.RecordNotFound {
				http.Error(rw, err.Error(), 404)
			} else {
				panic(err)
			}
		}

		if taskToComplete.CompletedDate.IsZero() {
			taskToComplete.CompletedDate = time.Now()
		} else {
			taskToComplete.CompletedDate = time.Time{}
		}

		err = db.Save(&taskToComplete).Error

		if err != nil {
			panic(err)
		}

		js, err := json.Marshal(taskToComplete)

		if err != nil {
			panic(err)
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(js)

	})
}
