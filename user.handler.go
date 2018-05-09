package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	service *UserDAO
)

func init() {
	service.Connect()
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
	err, existUser := service.findUserByEmail(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if &existUser != nil {
		fmt.Println("existing User: ", existUser)
		responseWithJson(w, http.StatusCreated, map[string]string{"message": "user existed"})
		return
	}

	pwd := []byte(user.Password)
	hassPassword, err := bcrypt.GenerateFromPassword(pwd, cost)
	user.Password = string(hassPassword)
	if er := service.insertUser(user); er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("user: ", user)
	responseWithJson(w, http.StatusCreated, map[string]string{"message": "succesful"})
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
	fmt.Println("email from request: ", email)
	err, user := service.findUserByEmail(email)
	if err != nil {
		// panic(err.Error())
		// Info.Println("error from find user by email: ", err.Error())
		// Info.Println("error from find user by email: ", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// log.Fatalln("user from find user: ", user)

	pwd := data["password"].(string)

	fmt.Println("password from request: ", pwd)

	isMatch := comparePasswords(user.Password, pwd)

	if isMatch == false {
		responseWithJson(w, http.StatusUnauthorized, map[string]string{"message": "incorrect passwod"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Email,
	})
	// log.Fatalln("token from request: ", token)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	fmt.Println("token string from request: ", tokenString)
	// json.NewEncoder(w).Encode(JwtToken{Token: tokenString})

	responseWithJson(w, http.StatusOK, map[string]string{"token": tokenString})

}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Print(w, "Hello There")
	err, users := service.findAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log.Fatalln("users: ", users)
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

	err = service.removeUser(&user)
	// data, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWithJson(w, http.StatusOK, map[string]string{"message": "removed"})
	return
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("email from request: ", email)
	// err, user := service.findUserByEmail(email)
	// if err != nil {
	// 	// panic(err.Error())
	// 	// Info.Println("error from find user by email: ", err.Error())
	// 	// Info.Println("error from find user by email: ", err.Error())
	// 	http.Error(w, err.Error(), http.StatusNotFound)
	// 	return
	// }
	// log.Fatalln("user from find user: ", user)

	pwd := data["password"].(string)

	bytepwd := []byte(pwd)
	hassPassword, err := bcrypt.GenerateFromPassword(bytepwd, cost)
	fmt.Println("password from request: ", pwd)

	err = service.updateUserPassword(email, string(hassPassword))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWithJson(w, http.StatusOK, map[string]string{"message": "change password completed"})
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// get email from request
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
	fmt.Println("email from request: ", email)
	// check if user exists
	err, user := service.findUserByEmail(email)
	if err != nil {
		// panic(err.Error())
		// Info.Println("error from find user by email: ", err.Error())
		// Info.Println("error from find user by email: ", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Println("user from find user: ", user)
	// if yes -> send email confirm with embbed encrypted link
	//
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	bytePwd := []byte(plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
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
