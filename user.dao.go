package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserDAO struct {
	Server   string
	Database string
}

type UserService interface {
	Connect()
	findAllUsers()
	insertUser(user User) error
	removeUser(user *User) error
	findUserByEmail(email string) (error, User)
}

var (
	db     *mgo.Database
	SERVER string

	instance *UserDAO
)

const (
	COLLECTION = "user"
	DATABASE   = "todoapp"
)

func init() {

	SERVER = os.Getenv("MLAB_URL")

	instance = &UserDAO{
		Server:   SERVER,
		Database: DATABASE,
	}

	// session, err := mgo.Dial(SERVER)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// db = session.DB(DATABASE)
}

func (userivce *UserDAO) Connect() {
	session, err := mgo.Dial(instance.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(instance.Database)
}

func (userivce *UserDAO) insertUser(user User) error {
	err := db.C(COLLECTION).Insert(&user)
	return err
}

func (userivce *UserDAO) findAllUsers() (error, []User) {
	var users []User
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	return err, users
}

func (userivce *UserDAO) removeUser(user *User) error {
	err := db.C(COLLECTION).Remove(user)
	return err
}

func (userivce *UserDAO) findUserByEmail(email string) (error, User) {
	// Info.Println("email passed ", email)
	var user User
	err := db.C(COLLECTION).Find(bson.M{"email": email}).One(&user)
	fmt.Println("Find user: ", user)
	return err, user
}
