package model

import (
	"test-CRUD/conn"

	"gopkg.in/mgo.v2/bson"
)

// User structure
type User struct {
	Id         string   `bson:"_id"`
	Password   string   `bson:"password"`
	IsActive   bool     `bson:"isActive"`
	Balance    string   `bson:"balance"`
	Age        int      `bson:"age"`
	Name       string   `bson:"name"`
	Gender     string   `bson:"gender"`
	Company    string   `bson:"company"`
	Email      string   `bson:"email"`
	Phone      string   `bson:"phone"`
	Address    string   `bson:"address"`
	About      string   `bson:"about"`
	Registered string   `bson:"registered"`
	Latitude   float64  `bson:"latitude"`
	Longitude  float64  `bson:"longitude"`
	Tags       []string `bson:"tags"`
	Friends    []Friend `bson:"friends"`
	Data       string   `bson:"data"`
}

type Friend struct {
	Id   int
	Name string
}

// Users list
type Users []User

// UserInfo model function
func UserInfo(id string, userCollection string) (User, error) {
	// Get DB from Mongo Config
	db := conn.GetMongoDB()
	user := User{}
	err := db.C(userCollection).Find(bson.M{"_id": &id}).One(&user)
	return user, err
}
