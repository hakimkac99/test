package user

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"test-CRUD/conn"
	user "test-CRUD/models/user"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// UserCollection statically declared
const UserCollection = "user"

var (
	errNotExist        = errors.New("Users are not exist")
	errInvalidID       = errors.New("Invalid ID")
	errInvalidBody     = errors.New("Invalid request body")
	errInsertionFailed = errors.New("Error in the user insertion")
	errUpdationFailed  = errors.New("Error in the user updation")
	errDeletionFailed  = errors.New("Error in the user deletion")
	errConexionFailed  = errors.New("Id or Password incorrect")
)

// InsertUsers Endpoint
/*
Here we choose to manage the deserialization of the age field which can be of type int or string.
We could also have chosen to return an error if the age field does not correspond to the type of the model. (int)
(same thing for the other fields of the user object)
*/
func InsertUsers(c *gin.Context) {
	//Variable for storing the input serialized JSON file
	var inputJSON []interface{}

	//deserialize JSON file into geniric interface (no types)
	err := c.BindJSON(&inputJSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInvalidBody.Error()})
		return
	}

	//db connection
	db := conn.GetMongoDB()

	//canal with length of the number of users which receive 1 if the user is inserted, 0 else.
	doneInsertingUser := make(chan int, len(inputJSON))

	for _, serialUser := range inputJSON {
		//for each user in the file, deserialize it in a specific Goroutine and verify the content before inserting in db.
		go func(serialUser map[string]interface{}) {

			//Variable to stock deserielized user
			deserialUser := user.User{}

			//Deserialize id field without controlling (validation)
			deserialUser.Id = serialUser["id"].(string)

			//check if the user exist in db
			u, _ := user.UserInfo(deserialUser.Id, UserCollection)

			//If user doesn't exist in db
			if u.Id == "" {

				//Hashing password in goroutine
				password := []byte(serialUser["password"].(string))
				doneHashing := make(chan int)
				go func() {
					hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
					if err != nil {
						panic(err)
					}
					deserialUser.Password = string(hashedPassword)
					doneHashing <- 1
				}()

				//Deserialize isActive and balance fields without controlling (validation)
				deserialUser.IsActive = serialUser["isActive"].(bool)
				deserialUser.Balance = serialUser["balance"].(string)

				//Deserialize age field with controlling
				switch serialUser["age"].(type) {
				case string:
					deserialUser.Age, _ = strconv.Atoi(serialUser["age"].(string))
				case float64:
					deserialUser.Age = int(serialUser["age"].(float64))
				}

				//Deserialize other fields without controlling
				deserialUser.Name = serialUser["name"].(string)
				deserialUser.Gender = serialUser["gender"].(string)
				deserialUser.Company = serialUser["company"].(string)
				deserialUser.Email = serialUser["email"].(string)
				deserialUser.Phone = serialUser["phone"].(string)
				deserialUser.Address = serialUser["address"].(string)
				deserialUser.About = serialUser["about"].(string)
				deserialUser.Registered = serialUser["registered"].(string)
				deserialUser.Latitude = serialUser["latitude"].(float64)
				deserialUser.Longitude = serialUser["longitude"].(float64)

				//Deserialize tags field
				for _, tag := range serialUser["tags"].([]interface{}) {
					deserialUser.Tags = append(deserialUser.Tags, fmt.Sprint(tag))
				}

				//Deserialize friends field
				for _, friend := range serialUser["friends"].([]interface{}) {
					friend := friend.(map[string]interface{})
					id := int(friend["id"].(float64))
					name := friend["name"].(string)
					deserialUser.Friends = append(deserialUser.Friends, user.Friend{id, name})
				}

				//Deserielize Data field without controlling
				deserialUser.Data = serialUser["data"].(string)

				//Goroutine for creating data in a file
				doneCreatingFile := make(chan int)
				go func() {
					fileName := "files/" + deserialUser.Id

					//Create the file in files folder (name of file is the id of the user)
					f, err := os.Create(fileName)
					if err != nil {
						panic(err)
					}
					defer f.Close()

					//write data to the file
					_, err = f.WriteString(deserialUser.Data)
					if err != nil {
						panic(err)
					}
					doneCreatingFile <- 1
				}()

				//Waiting the hashing to finish before inserting to database
				<-doneHashing

				//Inserting in database
				err = db.C(UserCollection).Insert(deserialUser)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInsertionFailed.Error()})
					return
				}

				//Waiting creating file before exiting the goroutine
				<-doneCreatingFile
				doneInsertingUser <- 1
			} else {
				//If user alredy exist
				doneInsertingUser <- 0
			}
		}(serialUser.(map[string]interface{}))

	}

	inserted, notInserted := 0, 0
	//wait for all goroutins to finish
	for i := 0; i < len(inputJSON); i++ {
		res := <-doneInsertingUser
		if res == 1 {
			inserted++
		} else {
			notInserted++
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "Success", "Inserted": inserted, "Not inserted": notInserted})

}

// GetAllUser Endpoint
func GetAllUsers(c *gin.Context) {
	// Get DB from Mongo Config
	db := conn.GetMongoDB()

	// Variables for storing all users
	users := user.Users{}

	//Request for retrieving data
	err := db.C(UserCollection).Find(bson.M{}).All(&users)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errNotExist.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "users": &users})
}

// GetUser Endpoint
func GetUser(c *gin.Context) {
	//get user info
	user, err := user.UserInfo(c.Param("id"), UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInvalidID.Error()})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": &user})
}

//Login Endpoint
//No JWT
func Login(c *gin.Context) {
	//struct for saving the deserialised JSON
	type form struct {
		Id       string
		Password string
	}
	loginForm := form{}

	//deserialize JSON form
	err := c.BindJSON(&loginForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInvalidBody.Error()})
		return
	}

	//Check if the user exist
	u, _ := user.UserInfo(loginForm.Id, UserCollection)
	if u.Id != "" {
		hashedPass := []byte(u.Password)
		formPass := []byte(loginForm.Password)

		//Compare the hashed password stored in db and password in the form
		err = bcrypt.CompareHashAndPassword(hashedPass, formPass)
		if err == nil {
			u.Password = ""
			c.JSON(http.StatusOK, gin.H{"status": "Connection success", "user": u})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errConexionFailed.Error()})
		}
	} else {
		//User doesn't exist
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errConexionFailed.Error()})
	}
}

// UpdateUser Endpoint
func UpdateUser(c *gin.Context) {
	// Get DB from Mongo Config
	db := conn.GetMongoDB()

	//Retrieve the existing user
	existingUser, err := user.UserInfo(c.Param("id"), UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInvalidID.Error()})
		return
	}
	oldData := existingUser.Data

	//Dserialize the JSON Form
	err = c.Bind(&existingUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInvalidBody.Error()})
		return
	}

	newData := existingUser.Data

	//Goroutine for updating file if data is changed
	doneUpdatingFile := make(chan int)
	go func() {
		if oldData != newData {
			fileName := "files/" + c.Param("id")
			e := os.Remove(fileName)
			if e != nil {
				log.Fatal(e)
			}

			f, err := os.Create(fileName)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			_, err = f.WriteString(newData)
			if err != nil {
				panic(err)
			}

			doneUpdatingFile <- 1
		} else {
			doneUpdatingFile <- 0
		}
	}()

	err = db.C(UserCollection).Update(bson.M{"_id": c.Param("id")}, existingUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errUpdationFailed.Error()})
		return
	}
	<-doneUpdatingFile
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": &existingUser})
}

// DeleteUser Endpoint
func DeleteUser(c *gin.Context) {
	// Get DB from Mongo Config
	db := conn.GetMongoDB()

	err := db.C(UserCollection).Remove(bson.M{"_id": c.Param("id")})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errDeletionFailed.Error()})
		return
	}

	//Delete file
	fileName := "files/" + c.Param("id")
	e := os.Remove(fileName)
	if e != nil {
		log.Fatal(e)
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User deleted successfully"})
}
