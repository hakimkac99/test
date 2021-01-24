package conn

import (
	"fmt"
	"os"

	mgo "gopkg.in/mgo.v2"
)

var db *mgo.Database

//MongoDB settings and connection
func init() {
	host := "localhost"
	dbName := "go-mongo"
	session, err := mgo.Dial(host)
	if err != nil {
		fmt.Println("session err:", err)
		os.Exit(2)
	}
	db = session.DB(dbName)
}

// GetMongoDB function to return DB connection
func GetMongoDB() *mgo.Database {
	return db
}
