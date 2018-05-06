package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	// "github.com/gin-gonic/gin"
	// _ "github.com/heroku/x/hmetrics/onload"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"

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

// var db *mgo.Database
// var ENCRYPT_KEY string

type ResponseParam struct {
	user    User   `json:"user"`
	users   []User `json:"users"`
	err     string `json:"error'`
	message string `json:"message"`
}

var (
	ENCRYPT_KEY string
	db          *mgo.Database
	SERVER      string

	Info *log.Logger
)

const (
	COLLECTION = "user"
	DATABASE   = "todoapp"
	cost       = 10
)

func Init(infoHandle io.Writer) {
	ENCRYPT_KEY = os.Getenv("ENCRYPT_KEY")
	// SERVER = os.Getenv("MLAB_URL")
	SERVER = "mongodb://duynguyen0428:cuongduy0428@ds221228.mlab.com:21228/todoapp"
	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	session, err := mgo.Dial(SERVER)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(DATABASE)
}

func main() {

	Init(os.Stdout)

	router := mux.NewRouter()
	router.Methods("GET").Path("/").HandlerFunc(IndexHandler)
	router.Methods("GET").Path("/user").HandlerFunc(GetAllUsersHandler)
	router.Methods("POST").Path("/user").HandlerFunc(CreatUserHandler)
	router.Methods("DELETE").Path("/user").HandlerFunc(RemoveUserHandler)
	router.Methods("POST").Path("/signin").HandlerFunc(SignInHandler)

	// http.HandleFunc("/", IndexHandler)
	// http.HandleFunc("/favicon.ico", FaviconHandler)
	// http.HandleFunc("/user", CreatUserHandler)
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, router)
	// http.ListenAndServe(":8000", router)

	// if port != "" {
	// 	http.ListenAndServe(":"+port, router)
	// } else {
	// 	http.ListenAndServe(":8000", router)
	// }

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
	pwd := []byte(user.Password)
	hassPassword, err := bcrypt.GenerateFromPassword(pwd, cost)
	user.Password = string(hassPassword)
	if er := insertUser(user); er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
		return
	}
	// response := ResponseParam{
	// 	message: "sucessfully",
	// }
	// response.message = "sucessfully"
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// data , err := json.NewEncoder(w).Encode(data)
	log.Fatalln("user: ", user)
	responseWithJson(w, http.StatusCreated, map[string]string{"message": "succesful"})
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Print(w, "Hello There")
	err, users := findAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Fatalln("users: ", users)
	responseWithJson(w, http.StatusOK, users)
	return
}

func RemoveUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = removeUser(&user)
	// data, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWithJson(w, http.StatusOK, map[string]string{"message": "removed"})
	return
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// unmarschal JSON
	err = json.Unmarshal(b, &data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := data["email"].(string)
	Info.Println("email from request: ", email)
	err, user := findUserByEmail(email)
	if err != nil {
		panic(err.Error())
		Info.Println("error from find user by email: ", err.Error())
		Info.Println("error from find user by email: ", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Fatalln("user from find user: ", user)

	pwd := []byte(data["password"].(string))
	log.Fatalln("password from request: ", pwd)

	isMatch := comparePasswords(user.Password, pwd)

	if isMatch == false {
		responseWithJson(w, http.StatusUnauthorized, map[string]string{"message": "incorrect passwod"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Email,
	})
	log.Fatalln("token from request: ", token)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		log.Fatalln(error)
	}
	log.Fatalln("token string from request: ", tokenString)
	// json.NewEncoder(w).Encode(JwtToken{Token: tokenString})

	responseWithJson(w, http.StatusOK, map[string]string{"token": tokenString})

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

func removeUser(user *User) error {
	err := db.C(COLLECTION).Remove(user)
	return err
}

func findUserByEmail(email string) (error, User) {
	log.Fatalln("email passed ", email)
	var user User
	err := db.C(COLLECTION).Find(bson.M{"Email": email}).One(&user)
	log.Fatalln("Find user: ", user)
	return err, user
}

// <=============== DAO ========================>

// <=============== Database ========================>

// <=============== Database ========================>

// <=============== Ultility Functions ========================>
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func responseWithJson(w http.ResponseWriter, reponsecode int, i interface{}) {
	response, _ := json.Marshal(i)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(i)
	w.Write(response)
	return
}
