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

// GetAllUser Endpoint
func GetAllUser(c *gin.Context) {
	// Get DB from Mongo Config
	db := conn.GetMongoDB()
	users := user.Users{}
	err := db.C(UserCollection).Find(bson.M{}).All(&users)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errNotExist.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "users": &users})
}

// GetUser Endpoint
func GetUser(c *gin.Context) {
	user, err := user.UserInfo(c.Param("id"), UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInvalidID.Error()})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"status": "success", "user": &user})
}

//Login Endpoint
func Login(c *gin.Context) {
	type form struct {
		Id       string
		Password string
	}
	loginForm := form{}

	err := c.BindJSON(&loginForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInvalidBody.Error()})
		return
	}

	u, _ := user.UserInfo(loginForm.Id, UserCollection)
	if u.Id != "" {
		hashedPass := []byte(u.Password)
		formPass := []byte(loginForm.Password)
		err = bcrypt.CompareHashAndPassword(hashedPass, formPass)
		if err == nil {
			u.Password = ""
			c.JSON(http.StatusOK, gin.H{"status": "Connexion success", "user": u})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errConexionFailed.Error()})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errConexionFailed.Error()})
	}
	fmt.Println(loginForm)
}

// InsertUsers Endpoint
func InsertUsers(c *gin.Context) {

	var inputJSON []interface{} //Variable for storing the input serialized JSON file

	err := c.BindJSON(&inputJSON) //deserialize JSON file into geniric interface (no types)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInvalidBody.Error()})
		return
	}

	db := conn.GetMongoDB()

	doneInsertingUser := make(chan int)

	for _, serialUser := range inputJSON {
		//for each user in the file => deserialize and verify the content before inserting in db in a specific Goroutine.
		go func(serialUser map[string]interface{}) {

			deserialUser := user.User{} //Variable to stock deserielized user

			deserialUser.Id = serialUser["id"].(string)

			//check if user exist
			u, _ := user.UserInfo(deserialUser.Id, UserCollection)

			if u.Id == "" { //User doesn't exist

				password := []byte(serialUser["password"].(string))

				doneHashing := make(chan int)

				//Hash password in goroutine
				go func() {
					hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
					if err != nil {
						panic(err)
					}
					deserialUser.Password = string(hashedPassword)
					doneHashing <- 1
				}()

				deserialUser.IsActive = serialUser["isActive"].(bool)
				deserialUser.Balance = serialUser["balance"].(string)

				//controle age (normalization)
				switch serialUser["age"].(type) {
				case string:
					deserialUser.Age, _ = strconv.Atoi(serialUser["age"].(string))
				case float64:
					deserialUser.Age = int(serialUser["age"].(float64))
				}

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

				for _, tag := range serialUser["tags"].([]interface{}) {
					deserialUser.Tags = append(deserialUser.Tags, fmt.Sprint(tag))
				}

				for _, friend := range serialUser["friends"].([]interface{}) {
					friend := friend.(map[string]interface{})
					id := int(friend["id"].(float64))
					name := friend["name"].(string)
					deserialUser.Friends = append(deserialUser.Friends, user.Friend{id, name})
				}

				deserialUser.Data = serialUser["data"].(string)

				doneCreatingFile := make(chan int)
				//Goroutine for creating data in file
				go func() {
					fileName := "files/" + deserialUser.Id

					f, err := os.Create(fileName)
					if err != nil {
						panic(err)
					}
					defer f.Close()
					_, err = f.WriteString(deserialUser.Data)
					if err != nil {
						panic(err)
					}
					doneCreatingFile <- 1
				}()

				<-doneHashing //Waiting the hashing to finish before inserting to database
				err = db.C(UserCollection).Insert(deserialUser)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInsertionFailed.Error()})
					return
				}

				//Waiting creating file before exiting the goroutine
				<-doneCreatingFile
				doneInsertingUser <- 1
			} else {
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

// UpdateUser Endpoint
func UpdateUser(c *gin.Context) {
	// Get DB from Mongo Config
	db := conn.GetMongoDB()

	existingUser, err := user.UserInfo(c.Param("id"), UserCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInvalidID.Error()})
		return
	}

	oldData := existingUser.Data

	err = c.Bind(&existingUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": errInvalidBody.Error()})
		return
	}

	newData := existingUser.Data

	doneUpdatingFile := make(chan int)
	if oldData != newData {
		go func() {
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
		}()
	}

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
	fileName := "files/" + c.Param("id")
	e := os.Remove(fileName)
	if e != nil {
		log.Fatal(e)
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User deleted successfully"})
}
