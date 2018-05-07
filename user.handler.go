package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

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
	Info.Println("email from request: ", email)
	err, user := findUserByEmail(email)
	if err != nil {
		panic(err.Error())
		Info.Println("error from find user by email: ", err.Error())
		Info.Println("error from find user by email: ", err.Error())
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
	err, users := findAllUsers()
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

	err = removeUser(&user)
	// data, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWithJson(w, http.StatusOK, map[string]string{"message": "removed"})
	return
}
