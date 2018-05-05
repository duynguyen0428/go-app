package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	// "github.com/gin-gonic/gin"
	// _ "github.com/heroku/x/hmetrics/onload"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

	router := mux.NewRouter()

	router.Methods("GET").Path("/user").HandlerFunc(GetAllUsersHandler)
	router.Methods("POST").Path("/user").HandlerFunc(CreatUserHandler)

	// http.HandleFunc("/", IndexHandler)
	// http.HandleFunc("/favicon.ico", FaviconHandler)
	// http.HandleFunc("/user", CreatUserHandler)

	http.ListenAndServe(":"+port, router)

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Print(w, "Hello There")
	err, users := findAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// data, err := json.Marshal(users)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	return
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

func findAllUsers() (error, []User) {
	var users []User
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	return err, users
}

// <=============== DAO ========================>

// <=============== Database ========================>

// <=============== Database ========================>
