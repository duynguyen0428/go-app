package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	// "github.com/gin-gonic/gin"
	// _ "github.com/heroku/x/hmetrics/onload"

	"gopkg.in/mgo.v2"
)

// <=============== Model ========================>
// User: user model
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDAO struct {
	Server   string
	Database string
}

// <=============== Model ========================>

var db *mgo.Database

const (
	COLLECTION = "user"
	SERVER     = "mongodb://duynguyen0428:cuongduy0428@ds221228.mlab.com:21228/todoapp"
	DATABASE   = "todoapp"
)

func init() {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(DATABASE)
}

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/favicon.ico", FaviconHandler)
	http.HandleFunc("/user", CreatUserHandler)

	http.ListenAndServe(":"+port, nil)

}

// <=============== Handlers ========================>
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Print(w, "Hello There")
	user := User{"test@mail.com", "123456"}
	data, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(data)
}
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./favicon.ico")
}

func CreatUserHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Print(w, "Hello There")

	if r.Method == "POST" {
		var user User
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if er := insertUser(user); er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	}
}

// <=============== Handlers ========================>

// <=============== DAO ========================>
// Establish a connection to database
// func (m *UserDAO) Connect() {
// 	session, err := mgo.Dial(m.Server)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	db = session.DB(m.Database)
// }
func insertUser(user User) error {
	err := db.C(COLLECTION).Insert(&user)
	return err
}

// <=============== DAO ========================>

// <=============== Database ========================>

// <=============== Database ========================>
